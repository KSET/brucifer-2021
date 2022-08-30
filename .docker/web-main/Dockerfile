# Statically build app with all needed dependencies
FROM node:18-alpine as asset-builder
WORKDIR /app
RUN apk update \
    && apk add \
      make
COPY assets .
RUN make assets

FROM golang:1.17 as builder

WORKDIR /go/src/app

RUN apt-get update \
    && apt-get install -y \
    make \
    upx

COPY go.mod go.sum ./

RUN CGO_ENABLED=0 go mod download

COPY . .
COPY --from=asset-builder assets .

RUN make build \
      OUTPUT_BINARY='/brucifer'


# Run app in a slim container
FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /brucifer ./brucifer

ENV PORT=3000

CMD ["./brucifer"]
