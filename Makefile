all: deps build test

deps:
	go get -v -d ./...

build:
	go build .

test:
	go test -v ./...

.PHONY: all build test
