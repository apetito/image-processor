FROM golang:1.16.4-alpine3.13 AS src

RUN apk update && apk upgrade && \
    apk add build-base && \
    go get -u "github.com/chai2010/webp" && \
    go get -u "github.com/anthonynsimon/bild"

WORKDIR /go/src/image_processor/

ENTRYPOINT ["/usr/local/go/bin/go", "run", "main.go"]
# --------------------------------------------------
FROM src AS build

COPY . /go/src/image_processor/

RUN go build -o /image_processor /go/src/image_processor/main.go
