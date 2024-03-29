bin = bin/bulletin
gofiles = $(shell find . -type f -name \*.go -or -name \*.gohtml)

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
		$(bin) -cache ./tmp/cache compose -template feed/page_template.gohtml -days 7 -output bulletin.tmp.html; \
	done

compose: $(bin)
	$(bin) --cache ./tmp/cache/ compose -days 7 -output bulletin.tmp.html

up: $(bin)
	zsh -eux -c "./update/update.sh > README.md 2> >(tee up.log)"
	git add README.md bulletins
	git commit -m "update"


up-push: up
	git push

uniq:
	sort feeds.conf | uniq > tmp && mv -fv tmp feeds.conf


.PHONY: clean smoke compose up up-push

