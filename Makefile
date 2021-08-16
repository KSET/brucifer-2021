OUTPUT_BINARY=bin/brucifer

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags osusergo,netgo -ldflags '-s -w -extldflags "-static"' -o "${OUTPUT_BINARY}" main.go

.PHONY: compact
compact: build
	upx --brute "${OUTPUT_BINARY}"
