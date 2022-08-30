LOCAL_UID=$(shell id -u)
LOCAL_GID=$(shell id -g)
OUTPUT_BINARY=bin/brucifer
DOCKER_COMPOSE=docker-compose --project-name 'brucifer-2021'
DOCKER_COMPOSE_DEV=$(DOCKER_COMPOSE) \
	-f 'docker-compose.yml' \
	-f 'docker-compose.dev.yml' \
	-f 'docker-compose.override.yml'

PACKAGE=brucosijada.kset.org
define LDFLAGS
-X '$(PACKAGE)/app/version.buildTimestamp=$(shell date --utc '+%Y-%m-%dT%H:%M:%S%z')'
endef
LDFLAGS:=$(strip $(LDFLAGS))

.PHONY: clean
clean: assets/clean
	rm -rf $(OUTPUT_BINARY)

.PHONY: build
build: assets
	CGO_ENABLED=0 \
	go \
	build \
	-a \
	-tags osusergo,netgo \
	-gcflags "all=-N -l" \
	-ldflags="-s -w -extldflags \"-static\" $(LDFLAGS)" \
	-o "${OUTPUT_BINARY}" \
	main.go

.PHONY: format
format:
	gofmt -e -l -s -w .

.PHONY: fmt
fmt: format

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
dev/server: assets
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
debug/build: assets
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
server/rebuild: docker/build server/start

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


# ASSETS
SHARP_CMD=npx sharp-cli
SHARP_FLAGS=--progressive --quality 100 --lossless --smartSubsample
SQUOOSH_CMD=npx @squoosh/cli
SQUOOSH_FLAGS=--optimizer-butteraugli-target 1.85
ASSET_LIST_BG=\
	assets/images/preview/bg.jpg
ASSET_LIST_BG_MOBILE=\
	assets/images/preview/bg-mobile.jpg
ASSET_LIST=$(ASSET_LIST_BG) $(ASSET_LIST_BG_MOBILE)

.PHONY: assets
assets: $(ASSET_LIST)

.PHONY: assets/clean
assets/clean:
	rm -rf $(ASSET_LIST)

$(ASSET_LIST_BG):
	$(SHARP_CMD) $(SHARP_FLAGS) \
		--input '$(strip $(patsubst %.jpg,,$(wildcard $(patsubst %.jpg,%,$@).*)))' \
		--output '$@' \
		resize 2560 \
	&& $(SQUOOSH_CMD) $(SQUOOSH_FLAGS) \
		--output-dir '$(@D)' \
		--resize '{"enabled":true,"width":2560,"height":1440,"method":"mitchell","fitMethod":"stretch","premultiply":true,"linearRGB":true}' \
		--mozjpeg '{"quality":85,"baseline":false,"arithmetic":false,"progressive":true,"optimize_coding":true,"smoothing":0,"color_space":3,"quant_table":3,"trellis_multipass":true,"trellis_opt_zero":true,"trellis_opt_table":true,"trellis_loops":1,"auto_subsample":true,"chroma_subsample":2,"separate_chroma_quality":false,"chroma_quality":75}' \
		$@

$(ASSET_LIST_BG_MOBILE):
	$(SHARP_CMD) $(SHARP_FLAGS) \
		--input '$(strip $(patsubst %.jpg,,$(wildcard $(patsubst %.jpg,%,$@).*)))' \
		--output '$@' \
		resize 1080 \
	&& $(SQUOOSH_CMD) $(SQUOOSH_FLAGS) \
		--output-dir '$(@D)' \
		--resize '{"enabled":true,"width":1080,"height":1920,"method":"mitchell","fitMethod":"stretch","premultiply":true,"linearRGB":true}' \
		--mozjpeg '{"quality":85,"baseline":false,"arithmetic":false,"progressive":true,"optimize_coding":true,"smoothing":0,"color_space":3,"quant_table":3,"trellis_multipass":true,"trellis_opt_zero":true,"trellis_opt_table":true,"trellis_loops":1,"auto_subsample":true,"chroma_subsample":2,"separate_chroma_quality":false,"chroma_quality":75}' \
		$@
