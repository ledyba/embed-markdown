.PHONY: all run get clean deploy

all:
	gofmt -w .
	go build -o embed-markdown github.com/ledyba/embed-markdown/...

run: all
	./embed-markdown

get:
	go get -u "github.com/russross/blackfriday"
	go get -u "github.com/microcosm-cc/bluemonday"

clean:
	go clean github.com/ledyba/embed-markdown/...

deploy:
	GOOS=linux GOARCH=amd64 go build -o embed-markdown github.com/ledyba/embed-markdown/...
	ssh ledyba.org mkdir -p /opt/run/embed-markdown
	scp embed-markdown embed-markdown.conf ledyba:/opt/run/embed-markdown
