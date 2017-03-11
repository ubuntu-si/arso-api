VERSION := 1.0.0
APP_NAME := arso

all: deps build

deps:
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/laher/gols/cmd/...
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	bin/gometalinter --install --update
	go get -t $(APP_NAME)/... # install test packages

sync:
	cd src/$(APP_NAME); glide install

update:
	cd src/$(APP_NAME); glide up

clean:
	rm -f $(APP_NAME)
	rm -rf pkg
	rm -rf bin
	find src/* -maxdepth 0 ! -name '$(APP_NAME)' -type d | xargs rm -rf
	rm -rf src/$(APP_NAME)/vendor/

build: sync
	go build --ldflags '-w -X main.build=$(VERSION)' $(APP_NAME)

lint:
	bin/gometalinter --fast --config=.golinter --cyclo-over=30 --deadline=60s --exclude $(shell pwd)/src/$(APP_NAME)/vendor src/$(APP_NAME)/...
	find src/$(APP_NAME) -not -path "./src/$(APP_NAME)/vendor/*" -name '*.go' | xargs gofmt -w -s

test: lint cover
	go test -v -race $(shell go-ls $(APP_NAME)/...)
	mv src/arso/main.apib src/arso/static/API.md

snowboard:
	wget https://github.com/subosito/snowboard/releases/download/vm1.0.0/snowboard-v1.0.0.linux-amd64.tar.gz -O - | tar -xz -C bin
	chmod +x bin/snowboard

docs:
	bin/snowboard lint src/arso/static/API.md
	bin/snowboard html -o src/arso/static/docs.html src/arso/static/API.md
	
cover:
	gocov test $(shell go-ls $(APP_NAME)/...) | gocov report