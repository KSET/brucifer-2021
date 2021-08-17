LOCAL_UID=$(shell id -u)
LOCAL_GID=$(shell id -g)
OUTPUT_BINARY=bin/brucifer

.PHONY: build
build:
	CGO_ENABLED=0 go build -a -tags osusergo,netgo -ldflags '-s -w -extldflags "-static"' -o "${OUTPUT_BINARY}" main.go

.PHONY: compact
compact: build
	upx --brute "${OUTPUT_BINARY}"

.PHONY: server/start
server/start: $docker/build $docker/up

.PHONY: server/stop
server/stop: $docker/down

.PHONY: $docker/build
$docker/build:
	docker-compose build \
    		--compress \
    		--parallel \
    		--pull

.PHONY: $docker/up
$docker/up:
	docker-compose \
		--project-name 'brucifer-2021' \
		up \
		--detach \
		--remove-orphans \
		--build

.PHONY: $docker/down
$docker/down:
	docker-compose \
		down \
		--remove-orphans
