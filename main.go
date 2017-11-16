package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"sync"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Item ...
type Item struct {
	url       string
	html      string
	updatedAt time.Time
}

// Cache ...
type Cache struct {
	entries map[string]*Item
	queue   []*Item
	mutex   *sync.Mutex
}

func newCache() *Cache {
	entries := make(map[string]*Item)
	queue := make([]*Item, 0)
	mutex := new(sync.Mutex)
	return &Cache{
		entries: entries,
		queue:   queue,
		mutex:   mutex,
	}
}

var cache *Cache
var lifeTime = flag.Int("lifetime", 20, "cache lifetime.")

func (cache *Cache) cleanUp(lifetime float64) {
	log.Info("Cache cleanup...")
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	idx := 0
	now := time.Now()
	total := len(cache.queue)
	for ; idx < len(cache.queue); idx++ {
		delta := now.Sub(cache.queue[idx].updatedAt)
		if delta.Minutes() < lifetime {
			break
		}
		delete(cache.entries, cache.queue[idx].url)
		log.Infof("Cache delete: %s", cache.queue[idx].url)
	}
	cache.queue = cache.queue[idx:]
	if idx > 0 {
		log.Infof("Delete: %d entries (from %d entries)", idx, total)
	}
	if len(cache.queue) != len(cache.entries) {
		log.Fatalf("Cache inconsistent: %d(queue) vs %d(entries)", len(cache.queue), len(cache.entries))
	}
}

func (cache *Cache) add(item *Item) {
	// Update cache
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.entries[item.url] = item
	cache.queue = append(cache.queue, item)
	log.Infof("Cache added: %s", item.url)
}

func (cache *Cache) find(url string) (*Item, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item, ok := cache.entries[url]
	return item, ok
}

func fetchURL(url string) (string, error) {
	if item, ok := cache.find(url); ok {
		return item.html, nil
	}
	log.Infof("Cache not found: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	unsafe := blackfriday.MarkdownCommon(body)
	html := string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	cache.add(&Item{
		url:       url,
		html:      html,
		updatedAt: time.Now(),
	})
	return html, nil
}

func encode(src string, async bool, id string) string {
	encoded, err := json.Marshal(src)
	body := ""
	if err != nil {
		log.Error("Oops.", err)
		body = err.Error()
	}
	body = string(encoded)
	switch async {
	case true:
		return fmt.Sprintf("document.getElementById('%s').innerHTML = (%s);", id, body)
	default:
		return fmt.Sprintf("document.write(%s);", body)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		index(w, r)
		return
	}
	paths := strings.Split(r.URL.Path, "/")
	async := false
	id := ""
	if len(paths) >= 2 && paths[len(paths)-2] == "async" {
		async = true
		id = paths[len(paths)-1]
	}
	url := r.URL.RawQuery
	if url == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "please specify url.")
		return
	}
	log.Infof("Path: %s / Query: %s", r.URL.Path, r.URL.RawQuery)
	log.Infof(" -> URL: %s", url)
	body, err := fetchURL(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, encode(err.Error(), async, id))
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprint(w, encode(body, async, id))
}

var port = flag.Int("port", 8080, "")

func startCacheDeleter() {
	c := make(chan bool)
	go func() {
		t := time.NewTicker(10 * time.Minute)
	END:
		for {
			select {
			case stop := <-c:
				if stop {
					break END
				}
			case <-t.C:
				cache.cleanUp(float64(*lifeTime))
			}
		}
		t.Stop()
	}()
}

func main() {
	flag.Parse() // Scan the arguments list
	cache = newCache()
	startCacheDeleter()
	http.Handle("/", http.StripPrefix("/", http.HandlerFunc(handler)))
	log.Printf("Start at http://localhost:%d/", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
