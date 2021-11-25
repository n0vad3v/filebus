# docker build . -t n0vad3v/filebus:latest
FROM golang:alpine as builder

RUN apk update && apk add alpine-sdk && mkdir /build
COPY go.mod /build
RUN cd /build && go mod download

COPY . /build
RUN cd /build && go build -o filebus .

FROM alpine

COPY --from=builder /build/filebus  /usr/bin/filebus

WORKDIR /data
CMD ["/usr/bin/filebus"]