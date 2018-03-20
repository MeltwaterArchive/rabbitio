VERSION := $(shell git describe --tags)
BUILD_DIR?=$(shell pwd)/build
NAME=rabbitio
DIRECTORIES=./ ./cmd ./rmq ./file

all: tools deps test

tools:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/mitchellh/gox

deps:
	dep ensure

test:
	go vet ${DIRECTORIES}
	go test ${DIRECTORIES}

build:
	go build -o ${NAME} -ldflags "-X main.version=${VERSION}" main.go


build-all:
	mkdir -p ${BUILD_DIR}/
	gox -verbose -ldflags "-X main.version=${VERSION}" \
	  -osarch="linux/amd64 darwin/amd64 windows/amd64 freebsd/amd64" \
	  -output="${BUILD_DIR}/${NAME}-${VERSION}-{{.OS}}-{{.Arch}}"

compress:
	gzip -f -v ${BUILD_DIR}/*

clean:
	rm -rf ./build
	rm -rf ./vendor

.PHONY: tools deps build-all compress build clean
