# FROM golang:1.22.0-alpine3.19 as build
# WORKDIR /go/src/app
# RUN apk update && apk add build-base
# # RUN useradd -u 909600094 tempuser
# # RUN groupmod -g 909600094 tempuser
# # RUN su tempuser
# COPY go.mod /go/src/app/
# COPY go.sum /go/src/app/
# RUN go mod download
# COPY . ./
# RUN go build
# FROM alpine
# WORKDIR /go/src/app
# COPY --from=build /go/src/app/api .
# COPY --from=build /go/src/app/database /go/src/app/database
# COPY --from=build /go/src/app/assets /go/src/app/assets

FROM --platform=linux/amd64 golang:1.22.0-alpine3.19