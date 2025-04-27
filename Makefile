SHELL:=/bin/bash

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

VERSION=dev
COMMIT=$(shell git rev-parse HEAD)
GITDIRTY=$(shell git diff --quiet || echo 'dirty')

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

DOCKER_REGISTRY := ""

APP_NAME := "spacetimedb-golang"
SERVICE_NAME := "spacetimedb-golang"
IMAGE := "spacetimedb-golang"

TARGET_MAX_CHAR_NUM=25
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

all: help

.PHONY: install-dependencies
## Install dependencies for the service
dependencies: ## Install dependencies for the service
	go mod tidy

.PHONY: build
## Build the binary for the service
build:
	CGO_ENABLED=1 go build -o ./bin/${APP_NAME} ./cmd/app/*.go

.PHONY: run-example
## Build and run the service binary
run-example: build
	APP_ENV=dev ./bin/${APP_NAME}

.PHONY: test
## Run the tests
test:
	go test -v ./...

.PHONY: lint
## Run the tests
lint:
	golangci-lint run ./...
