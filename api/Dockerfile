FROM golang:1.22.0-alpine3.19 as build
WORKDIR /go/src/app
RUN apk update && apk add build-base
COPY go.mod go.sum /go/src/app/
RUN go mod download
COPY . ./
RUN go build
FROM alpine
WORKDIR /go/src/app
COPY --from=build /go/src/app/api .
COPY --from=build /go/src/app/database /go/src/app/database
COPY --from=build /go/src/app/assets /go/src/app/assets

