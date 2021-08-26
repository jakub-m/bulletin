bin = bin/bulletin
gofiles = $(shell find . -type f -name \*.go)

default: test build

build: $(bin)

test: $(gofiles)
	go test ./...

fmt: $(gofiles)
	go fmt ./...

$(bin): $(gofiles)
	mkdir -p bin
	go build -o $(bin) cli/bulletin/bulletin.go

clean:
	rm -rf bin tmp

smoke: clean build fetch compose

fetch:
	mkdir -p tmp/cache
	$(bin) --cache ./tmp/cache/ fetch --url http://googleaiblog.blogspot.com/atom.xml
	$(bin) --cache ./tmp/cache/ fetch --url https://netflixtechblog.com/feed
	$(bin) --cache ./tmp/cache/ fetch --url http://muratbuffalo.blogspot.com/feeds/posts/default
	$(bin) --cache ./tmp/cache/ fetch --url https://perspectives.mvdirona.com/feed/
	$(bin) --cache ./tmp/cache/ fetch --url https://berthub.eu/articles/index.xml

compose: $(bin)
	$(bin) --cache ./tmp/cache/ compose --days 31 | tee bulletin.tmp.html

.PHONY: clean smoke fetch compose

