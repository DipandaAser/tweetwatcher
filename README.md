# tweetwatcher

A simple Telegram bot that send you a screenshot when there is a new tweet with some specified hashtag.

## How to run

First create a telegram bot. See the guide [here](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
Second create the .env file using env.example as template

### With binaries

You can grab binaries in the [releases](https://github.com/DipandaAser/tweetwatcher/releases) section.

Put the .env file in the same directory with the binary and run the binary

### With sources

Alternatively, to get latest and greatest run:

`go get -u github.com/DipandaAser/tweetwatcher`

- Prerequisites
    - Golang 1.15+
    - MongoDB

- Get dependencies
```shell
go get -d -v ./...
```

- Build
```shell
go build - o app
```

* Run
```shell
./app
```
