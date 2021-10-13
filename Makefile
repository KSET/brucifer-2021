LOCAL_UID=$(shell id -u)
LOCAL_GID=$(shell id -g)
OUTPUT_BINARY=bin/brucifer
DOCKER_COMPOSE=docker-compose --project-name 'brucifer-2021'
DOCKER_COMPOSE_DEV=$(DOCKER_COMPOSE) \
	-f 'docker-compose.yml' \
	-f 'docker-compose.dev.yml' \
	-f 'docker-compose.override.yml'

.PHONY: build
build:
	CGO_ENABLED=0 \
	go \
	build \
	-a \
	-tags osusergo,netgo \
	-gcflags "all=-N -l" \
	-ldflags '-s -w -extldflags "-static"' \
	-o "${OUTPUT_BINARY}" \
	main.go

.PHONY: compact
compact: build
	upx --brute "${OUTPUT_BINARY}"

.PHONY: sync-deps
sync-deps:
	CGO_ENABLED=0 go mod download

.PHONY: $pull
$pull:
	git pull --rebase

.PHONY: dev/server/start
dev/server/start:
	$(DOCKER_COMPOSE_DEV) \
		up \
		--detach \
		--remove-orphans \
		--build \
	&& \
	$(DOCKER_COMPOSE_DEV) \
		logs \
		-f \
		'web-main' \
		'web-admin' \
	|| \
	$(MAKE) dev/server/stop

.PHONY: dev/server/stop
dev/server/stop:
	$(DOCKER_COMPOSE_DEV) \
    		down \
    		--timeout 3 \
    		--remove-orphans

.PHONY: dev/server
dev/server:
	gin \
		--all \
		--immediate \
		--excludeDir 'admin-ui' \
		--excludeDir '.idea' \
		--excludeDir '.cache' \
		--excludeDir '.docker' \
		--bin "${OUTPUT_BINARY}" \
		run \
		main.go

.PHONY: debug/build
debug/build:
	CGO_ENABLED=0 \
	go \
	build \
	-a \
	-gcflags "all=-N -l" \
	-o "${OUTPUT_BINARY}" \
	main.go

.PHONY: debug/run
debug/run:
	killall dlv; \
	dlv \
	--listen=:2345 \
	--headless=true \
	--api-version=2 \
	--accept-multiclient \
	exec \
	"${OUTPUT_BINARY}"

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
	$(DOCKER_COMPOSE) \
			build \
    		--compress \
    		--pull

.PHONY: docker/up
docker/up:
	$(DOCKER_COMPOSE) \
		up \
		--detach \
		--remove-orphans

.PHONY: docker/down
docker/down:
	$(DOCKER_COMPOSE) \
		down \
		--remove-orphans
