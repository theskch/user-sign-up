# This version-strategy uses git tags to set the version string
# VERSION := $(shell git describe --tags --always --dirty)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

SWAGGER_YML=index.yml
SWAGGER_UI_DIST=rest/swagger-ui-dist

start: ## Start server.
	go run main.go
binary: ## [SWAGGER] Generate swagger binary.
	cp rest/api-interface/$(SWAGGER_YML) ./$(SWAGGER_UI_DIST)
	statik -src=$(SWAGGER_UI_DIST)/ -p=swaggerui
	mv swaggerui/statik.go rest/swaggerui/
	rm -r swaggerui
validate: ## [SWAGGER] Validate swagger
	swagger validate rest/api-interface/$(SWAGGER_YML)
generate-model: ## [SWAGGER] Generate swagger binary
	swagger generate model -f=./rest/api-interface/$(SWAGGER_YML)
test: ## Unit tests.
	go test -v ./...

build: validate binary generate-model