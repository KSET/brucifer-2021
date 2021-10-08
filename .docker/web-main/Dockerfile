# Statically build app with all needed dependencies
FROM golang:1.17 as builder

WORKDIR /go/src/app

RUN apt update \
    && apt install \
    upx \
    make

COPY Makefile go.mod go.sum ./

RUN make sync-deps

COPY . .

RUN make build


# Run app in a slim container
FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /go/src/app/bin/brucifer ./

ENV PORT=3000

CMD ["./brucifer"]
