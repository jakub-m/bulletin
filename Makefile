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
	echo "http://googleaiblog.blogspot.com/atom.xml\n\
https://netflixtechblog.com/feed\n\
http://muratbuffalo.blogspot.com/feeds/posts/default\n\
https://perspectives.mvdirona.com/feed/" | \
	$(bin) --cache ./tmp/cache/ --verbose fetch --feeds -

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

update-bulletin: $(bin)
	(cd bulletins; ../bin/bulletin)


.PHONY: clean smoke compose

