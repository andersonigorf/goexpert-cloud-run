IMAGE_NAME ?= cloudrun-goexpert:latest

run:
	@echo "Starting Docker container..."
	docker-compose up -d

run-prod:
	@echo "Starting Docker container..."
	docker-compose -f docker-compose.prod.yaml up -d

build:
	@echo "Building docker image $(IMAGE_NAME)..."
	docker build -t $(IMAGE_NAME) -f Dockerfile.prod .

test:
	go test ./... -v

.PHONY: run run-prod build test