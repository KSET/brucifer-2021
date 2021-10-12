# Statically build app with all needed dependencies
FROM golang:1.17 as builder

WORKDIR /go/src/app

RUN apt-get update \
    && apt-get install -y \
    upx

COPY go.mod go.sum ./

RUN CGO_ENABLED=0 go mod download

COPY . .

RUN CGO_ENABLED=0 \
    	go \
    	build \
    	-a \
    	-gcflags "all=-N -l" \
    	-o "/brucifer" \
    	main.go


# Run app in a slim container
FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /brucifer ./brucifer

ENV PORT=3000

CMD ["./brucifer"]
