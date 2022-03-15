.PHONY: run
run:
	make up_postgres
	make test
	make build
	docker build -t shortener .
	docker run shortener

.PHONY: build
build:
	go build -o "shortener" cmd/shortener/main.go

.PHONY: test
test:
	make up_postgres
	go test -v -race -timeout 30.0s ./...

.PHONY: up_postgres
up_postgres:
	./postgres_starter.sh

.DEFAULT_GOAL := run
