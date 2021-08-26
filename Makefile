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

fetch: $(bin)
	mkdir -p tmp/cache
	$(bin) --cache ./tmp/cache/ fetch -- \
		http://googleaiblog.blogspot.com/atom.xml \
		https://netflixtechblog.com/feed \
		http://muratbuffalo.blogspot.com/feeds/posts/default \
		https://perspectives.mvdirona.com/feed/ \
		https://berthub.eu/articles/index.xml \
		https://josephg.com/blog/rss/ \
		https://blog.khanacademy.org/engineering/rss \
		https://stackoverflow.blog/feed/ \
		https://tailscale.com/blog/index.xml \
		https://earthly.dev/blog/feed.xml \
		https://medium.com/@cep21/feed \
		https://dropbox.tech/feed \
		https://logicai.io/blog/feed \
		https://buttondown.email/hillelwayne/rss


compose: $(bin)
	$(bin) --cache ./tmp/cache/ compose --days 31 | tee bulletin.tmp.html

.PHONY: clean smoke compose

