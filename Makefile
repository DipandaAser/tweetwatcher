APPNAME=tweetwatcher
.DEFAULT_GOAL := help

run:
	go run main.go

## build: build application binary.
build:build-windows-32 build-windows-64 build-linux-32 build-linux-64 build-macOS-64
	@echo "Build done!"

build-windows-64:
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -o bin/$(APPNAME)-windows-amd64.exe

build-windows-32:
	go env -w GOOS=windows
	go env -w GOARCH=386
	go build -o bin/$(APPNAME)-windows-386.exe


build-macOS-64:
	go env -w GOOS=darwin
	go env -w GOARCH=amd64
	go build -o bin/$(APPNAME)-darwin-amd64

build-linux-64:
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o bin/$(APPNAME)-linux-amd64

build-linux-32:
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o bin/$(APPNAME)-linux-386

## docker-build: build the docker image
.PHONY: docker-build
docker-build:
	docker build -t tweetwatcher .

## docker-run: run the docker container in foreground
.PHONY: docker-run
docker-run:
	docker run -it --rm --name tweetwatcher --env-file .env tweetwatcher

## docker-build: run the docker container in background
.PHONY: docker-run-d
docker-run-d:
	docker run -d --rm --name tweetwatcher --env-file .env tweetwatcher

all: help
.PHONY: help
help: Makefile
	@echo " Choose a command..."
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'