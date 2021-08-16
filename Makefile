OUTPUT_BINARY=bin/brucifer

.PHONY: build
build:
	go build -ldflags '-s -w' -o "${OUTPUT_BINARY}" main.go

.PHONY: build
compact: build
	upx --brute "${OUTPUT_BINARY}"
