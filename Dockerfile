FROM golang:1.21.10-alpine AS build

WORKDIR /app

COPY go.* /app/
# If you are in China, you can use the following command to speed up the download
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . /app
RUN go build -o example

FROM alpine:latest

RUN mkdir /app
COPY --from=build /app/example /app
COPY --from=build /app/tls/ /app/tls/
