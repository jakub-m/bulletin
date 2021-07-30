default: test build
test:
	go test ./...
build:
	mkdir -p bin
	#go build cli/fetchurl/fetchurl.go
	go build -o bin/feedsummary cli/feedsummary/feedsummary.go
fmt:
	go fmt ./...
clean:
	rm -rf bin
smoke: test build
	bin/feedsummary fetch --url http://googleaiblog.blogspot.com/atom.xml


.PHONY: test build

