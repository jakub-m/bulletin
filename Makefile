default: test build
test:
	go test ./...
build:
	go build cli/fetchurl/fetchurl.go

.PHONY: test build

