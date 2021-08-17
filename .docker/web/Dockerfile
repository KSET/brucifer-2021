# Statically build app with all needed dependencies
FROM golang:1.17-alpine as builder

WORKDIR /go/src/app

RUN apk update \
    && apk add --update \
    upx \
    make

COPY . .

RUN make build


# Run app in a slim container
FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /go/src/app/bin/brucifer ./

CMD ["./brucifer"]
