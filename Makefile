.PHONY: build
build:
	go build -v ./cmd/shortener

.PHONY: test
test: 
	go test -v -race -timeout 30.0s ./...

.DEFAULT_GOAL := build
