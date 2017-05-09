.PHONY: all run get clean

all:
	gofmt -w .
	go build -o embed-markdown github.com/ledyba/embed-markdown/...

run: all
	./embed-markdown

get:
	go get -u "github.com/russross/blackfriday"
	go get -u "github.com/microcosm-cc/bluemonday"
	go get -u "github.com/Sirupsen/logrus"

clean:
	go clean github.com/ledyba/embed-markdown/...
