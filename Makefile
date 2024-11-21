DOCKER_REGISTRY := registry.wantia.app
DOCKER_IMAGE := mwantia/nautilus
DOCKER_VERSION := v0.0.1
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64

.PHONY: all setup test release cleanup

all: cleanup setup release

setup:
	docker buildx create --use --name multi-arch-builder || true

test:
	go run cmd/nautilus/main.go --config tests/nautilus.yml

release: setup
	docker buildx build --push --platform ${DOCKER_PLATFORMS} -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_VERSION) -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest .

cleanup:
	rm -rf tmp/*