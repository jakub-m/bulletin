default: test build
test:
	go test ./...
build:
	mkdir -p bin
	go build -o bin/bulletin cli/bulletin/bulletin.go
fmt:
	go fmt ./...
clean:
	rm -rf bin tmp
smoke: clean test build
	mkdir -p tmp/cache
	#bin/bulletin --cache ./tmp/cache/ fetch --url http://googleaiblog.blogspot.com/atom.xml
	bin/bulletin --cache ./tmp/cache/ fetch --url https://netflixtechblog.com/feed
	bin/bulletin --cache ./tmp/cache/ compose

.PHONY: test build

