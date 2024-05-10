# syntax=docker/dockerfile:1
FROM --platform=linux/amd64 golang:alpine AS builder

WORKDIR /app

RUN apk update && \
    apk upgrade && \
    apk add gcc musl-dev && \
    apk add sqlite-dev vips-dev && \
    apk cache clean

RUN mkdir -p /app

COPY app /app/app
COPY assets /app/assets
COPY go.mod /app/go.mod
COPY go.work /app/go.word
COPY imgs /app/imgs
COPY main.go /app/main.go
COPY pkg /app/pkg
COPY templates /app/templates

WORKDIR /app

RUN go get . && \
    go build -o bin/start . 

CMD ["/app/bin/start"]
