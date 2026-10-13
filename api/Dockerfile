FROM golang:1.27-alpine as build
WORKDIR /jovvix
RUN apk update && apk add build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o jovvix
FROM alpine
WORKDIR /jovvix
COPY --from=build /jovvix/jovvix .
COPY --from=build /jovvix/database ./database
COPY --from=build /jovvix/assets ./assets
