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

smoke: clean build fetch-smoke compose

fetch-smoke: $(bin)
	mkdir -p tmp/cache
	$(bin) --cache ./tmp/cache/ --verbose fetch -- \
		http://googleaiblog.blogspot.com/atom.xml \
		https://netflixtechblog.com/feed \
		http://muratbuffalo.blogspot.com/feeds/posts/default \
		https://perspectives.mvdirona.com/feed/ \

fetch: $(bin)
	mkdir -p tmp/cache
	$(bin) --cache ./tmp/cache/ --verbose fetch  --feeds feeds.conf

watch-template: $(bin)
	while [ 1 ]; do 
		echo watch...; \
		sleep 1; \
		fswatch -1 feed/page_template.gohtml && \
		$(bin) -cache ./tmp/cache compose -template feed/page_template.gohtml -days 7 | tee bulletin.tmp.html; \
	done

compose: $(bin)
	$(bin) --cache ./tmp/cache/ compose --days 7 | tee bulletin.tmp.html

.PHONY: clean smoke compose

