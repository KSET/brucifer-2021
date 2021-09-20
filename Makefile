LOCAL_UID=$(shell id -u)
LOCAL_GID=$(shell id -g)
OUTPUT_BINARY=bin/brucifer

.PHONY: build
build:
	CGO_ENABLED=0 go build -a -tags osusergo,netgo -ldflags '-s -w -extldflags "-static"' -o "${OUTPUT_BINARY}" main.go

.PHONY: compact
compact: build
	upx --brute "${OUTPUT_BINARY}"

.PHONY: $pull
$pull:
	git pull --rebase

.PHONY: dev/server
dev/server:
	gin \
		--all \
		--immediate \
		--bin "${OUTPUT_BINARY}" \
		run \
		main.go

.PHONY: server/start
server/start: docker/up

.PHONY: server/stop
server/stop: docker/down

.PHONY: server/restart
server/restart: server/stop server/start

.PHONY: server/rebuild
server/rebuild: docker/build server/restart

.PHONY: server/redeploy
server/redeploy: $pull server/rebuild

.PHONY: docker/build
docker/build:
	docker-compose build \
    		--compress \
    		--pull

.PHONY: docker/up
docker/up:
	docker-compose \
		--project-name 'brucifer-2021' \
		up \
		--detach \
		--remove-orphans \
		--build

.PHONY: docker/down
docker/down:
	docker-compose \
		down \
		--remove-orphans
