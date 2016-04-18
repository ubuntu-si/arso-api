VERSION := 1.0.0

all: setup build

wercker: all test

setup:
	go get github.com/axw/gocov/gocov
	go get github.com/smartystreets/goconvey
	go get -v arso

clean:
	rm -f arso
	rm -rf pkg
	rm -rf bin
	find src/* -maxdepth 0 ! -name 'arso' -type d | xargs rm -rf

build:
	go build --ldflags '-w -X main.build=$(VERSION)+$(shell git rev-parse --short HEAD)' arso

test:
	go test -v -race arso/...

cover:
	gocov test arso/... | gocov report
