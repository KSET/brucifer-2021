# Statically build app with all needed dependencies
FROM node:16 as asset-builder
RUN apt-get update && apt-get install -y libvips-dev --no-install-recommends
WORKDIR /app
RUN npm install --global \
        sharp-cli \
        @squoosh/cli
COPY assets ./assets
COPY ./Makefile .
RUN make -j assets \
      SHARP_CMD=sharp \
      SQUOOSH_CMD=squoosh-cli

FROM golang:1.19 as builder
WORKDIR /go/src/app
RUN apt-get update \
    && apt-get install -y \
    make \
    upx
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 go mod download
COPY . .
COPY --from=asset-builder /app/assets ./assets
RUN make build \
      OUTPUT_BINARY='/brucifer'


# Run app in a slim container
FROM alpine:latest as runner
WORKDIR /root/
COPY --from=builder /brucifer ./brucifer
ENV PORT=3000
CMD ["./brucifer"]
