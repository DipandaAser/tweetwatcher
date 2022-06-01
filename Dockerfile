FROM golang:alpine AS base

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

#copies go.sum and go.mod if exists
COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o tweetwatcher .
RUN apk add --no-cache ca-certificates

# build image with the binary
FROM scratch

# copy certificate to be able to make https request to twitter and telegram
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /build/tweetwatcher /

ENTRYPOINT ["/tweetwatcher"]