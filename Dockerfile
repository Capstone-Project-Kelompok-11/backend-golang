# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV DOCKERIZED 1

COPY .env.docker /
COPY .env.docker /app
RUN mkdir -p /app/bin
COPY bin/start /app/bin

CMD ["/app/bin/start"]

FROM scratch

ENV DOCKERIZED 1

COPY --from=builder /app/bin/start /usr/local/bin/start

CMD ["/usr/local/bin/start"]