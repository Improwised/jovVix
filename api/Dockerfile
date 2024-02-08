# syntax = docker/dockerfile:latest
FROM golang:1.22.0-alpine3.19 AS build

ARG GOPATH="/go"
ARG GOMODCACHE=${GOPATH}/pkg/mod
ARG GOCACHE="/root/.cache/go-build"

ENV GOMODCACHE=${GOMODCACHE}
ENV GOCACHE=${GOCACHE}

WORKDIR /go/src/app

RUN apk update && \
    apk add build-base

COPY go.* ./

RUN --mount=type=cache,mode=0777,target=${GOCACHE} \
    --mount=type=cache,mode=0777,target=${GOMODCACHE} \
    go mod download

COPY . ./

RUN --mount=type=cache,mode=0777,target=${GOCACHE} \
    --mount=type=cache,mode=0777,target=${GOMODCACHE} \
    go build -o /tmp/app

FROM alpine
WORKDIR /app

ENV MODE="docker"
ENV DOCKERIZE_VERSION=v0.6.1

RUN set -ex; \
    apk update && \
    apk add --no-cache wget build-base && \
    wget https://github.com/jwilder/dockerize/releases/download/${DOCKERIZE_VERSION}/dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz && \
    tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz && \
    rm dockerize-alpine-linux-amd64-${DOCKERIZE_VERSION}.tar.gz

COPY assets database .env.example ./
COPY --from=build /tmp/app ./

ENTRYPOINT [ "dockerize", "-template", ".env.example:.env", "./app", "api"]

CMD ["./app", "api"]
