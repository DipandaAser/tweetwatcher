APPNAME=tweetwatcher

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