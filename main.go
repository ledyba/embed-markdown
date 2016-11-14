package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
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

var cache map[string]*Item
var queue []*Item
var cacheMutex *sync.Mutex

func cacheCleanUp() {
	log.Info("Cache cleanup...")
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	idx := 0
	now := time.Now()
	total := len(queue)
	for ; idx < len(queue); idx++ {
		delta := now.Sub(queue[idx].updatedAt)
		if delta.Minutes() < 20 {
			break
		}
		delete(cache, queue[idx].url)
		log.Infof("Cache delete: %s", queue[idx].url)
	}
	queue = queue[idx:]
	if idx > 0 {
		log.Infof("Delete: %d entries (from %d entries)", idx, total)
	}
	if len(queue) != len(cache) {
		log.Fatalf("Cache inconsistent: %d(queue) vs %d(cache)", len(queue), len(cache))
	}
}

func cacheAdd(item *Item) {
	// Update cache
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache[item.url] = item
	queue = append(queue, item)
	log.Info("Cache added: %s", item.url)
}

func fetchURL(url string) (string, error) {
	if item, ok := cache[url]; ok {
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
	cacheAdd(&Item{
		url:       url,
		html:      html,
		updatedAt: time.Now(),
	})
	return html, nil
}

func encode(src string) string {
	encoded, err := json.Marshal(src)
	body := ""
	if err != nil {
		log.Error("Oops.", err)
		body = err.Error()
	}
	body = string(encoded)
	return fmt.Sprintf("document.write(%s);", body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RawQuery
	if url == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "please specify url.")
		return
	}
	log.Infof("URL: %s", url)
	body, err := fetchURL(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, encode(err.Error()))
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, encode(body))
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
				cacheCleanUp()
			}
		}
		t.Stop()
	}()
}

func main() {
	flag.Parse() // Scan the arguments list
	cache = make(map[string]*Item)
	queue = make([]*Item, 0)
	cacheMutex = new(sync.Mutex)
	startCacheDeleter()
	http.HandleFunc("/", handler)
	log.Printf("Start at http://localhost:%d/", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
