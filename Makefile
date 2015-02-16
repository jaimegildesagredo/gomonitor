all: build test

build:
	go build .

test:
	go test -v ./...

.PHONY: all build test
