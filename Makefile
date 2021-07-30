default: test build
test:
	go test ./...
build:
	mkdir -p bin
	go build -o bin/feedsummary cli/feedsummary/feedsummary.go
fmt:
	go fmt ./...
clean:
	rm -rf bin tmp_cache
smoke: clean test build
	rm -rf tmp_cache
	mkdir -p tmp_cache 
	bin/feedsummary --cache ./tmp_cache/ fetch --url http://googleaiblog.blogspot.com/atom.xml


.PHONY: test build

