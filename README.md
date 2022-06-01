# tweetwatcher

A simple Telegram bot that send you a screenshot when there is a new tweet with some specified hashtag.

## How to run

First create a telegram bot. See the guide [here](https://core.telegram.org/bots#3-how-do-i-create-a-bot)

Second create the .env file using env.example as template

### With binaries

You can grab binaries in the [releases](https://github.com/DipandaAser/tweetwatcher/releases) section.

Put the .env file in the same directory with the binary and run the binary

### With sources

- Prerequisites
  - Golang 1.15+

- Clone the repo
```shell
git clone https://github.com/DipandaAser/tweetwatcher.git
cd  tweetwatcher
```

- Get dependencies
```shell
go mod download
```

- Run
```shell
go run main.go
```

### With Docker

- Build the image
```shell
make docker-build
```

- Run the container
  - In foreground
  ```shell
  make docker-run
  ```
    - In background
  ```shell
  make docker-run-d
  ```