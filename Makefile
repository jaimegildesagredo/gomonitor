all: deps build test

deps:
	go get -v -d ./...

build:
	sed -i '/var NETWORK_DASHBOARD_HTML/d' gomonitor.go
	echo "var NETWORK_DASHBOARD_HTML = `go run utils/f2g.go dashboards/network/index.html | sed 's/\ /\,\ /g' | sed 's/^\[/\[]byte{/' | sed 's/\]$$/\}/'`" >> gomonitor.go
	go build .

test:
	go test -v ./...

.PHONY: all build test
