DOCKER_REGISTRY := registry.wantia.app
DOCKER_IMAGE := mwantia/nautilus
DOCKER_VERSION := v0.0.2
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64
BUILD_NAME := nautilus

.PHONY: build test docker-setup docker-release docker-cleanup

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o build/$(BUILD_NAME) ./cmd/nautilus/main.go ;

test: build
	./build/nautilus agent --config ./tests/config.hcl

plugin-debug: build
	./build/nautilus plugin debug --address http://127.0.0.1:8080

docker-setup:
	docker buildx create --use --name multi-arch-builder || true

docker-build: docker-setup
	docker buildx build --platform ${DOCKER_PLATFORMS} -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_VERSION) -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest .

docker-release: docker-setup
	docker buildx build --push --platform ${DOCKER_PLATFORMS} -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_VERSION) -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest .