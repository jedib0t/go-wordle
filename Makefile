.PHONY: all test

default: run

all: test

build:
	go build .

fmt:
	go fmt $(shell go list ./...)

lint:
	golint -set_exit_status $(shell go list ./...)

run:
	go run .

test: fmt lint vet build
	go test -cover -coverprofile=.coverprofile $(shell go list ./...)

vet:
	go vet $(shell go list ./...)

