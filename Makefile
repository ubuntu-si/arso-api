all: deps build

wercker: all test

deps:
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/laher/gols/cmd/...
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	bin/gometalinter --install --update
	go get -t arso/... # install test packages

sync:
	cd src/arso; govendor sync

vendor:
	cd src/arso; govendor update +external
	git add src/arso/vendor/vendor.json -f

clean:
	rm -f arso
	rm -rf pkg
	rm -rf bin
	find src/* -maxdepth 0 ! -name 'arso' -type d | xargs rm -rf
	find src/arso/vendor/* -maxdepth 0 ! -name 'arso' -type d | xargs rm -rf

build: sync
	go build --ldflags '-w' arso

lint:
	bin/gometalinter --fast --disable=gotype --disable=gas --disable=dupl --cyclo-over=30 --deadline=60s --exclude $(shell pwd)/src/arso/vendor src/arso/...
	find src/arso -not -path "./src/arso/vendor/*" -name '*.go' | xargs gofmt -w -s

test: lint cover
	go test -v -race $(shell go-ls arso/...)

cover:
	gocov test $(shell go-ls arso/...) | gocov report

editor:
	go get -u -v github.com/nsf/gocode
	go get -u -v github.com/rogpeppe/godef
	go get -u -v github.com/golang/lint/golint
	go get -u -v github.com/lukehoban/go-outline
	go get -u -v sourcegraph.com/sqs/goreturns
	go get -u -v golang.org/x/tools/cmd/gorename
	go get -u -v github.com/tpng/gopkgs
	go get -u -v github.com/newhook/go-symbols
	go get -u -v golang.org/x/tools/cmd/guru
