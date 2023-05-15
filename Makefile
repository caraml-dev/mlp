include .env.development
export

UI_PATH := ui
UI_BUILD_PATH := ${UI_PATH}/build
API_PATH := api
API_ALL_PACKAGES := $(shell cd ${API_PATH} && go list ./... | grep -v github.com/caraml-dev/mlp/client | grep -v mocks)
BIN_NAME := $(if ${APP_NAME},${APP_NAME},mlp)

all: setup init-dep lint test clean build run

# ============================================================
# Initialize dependency recipes
# ============================================================
.PHONY: setup
setup:
	@echo "> Setting up tools..."
	@test -x $(shell go env GOPATH)/bin/golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.48.0/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48.0

.PHONY: init-dep
init-dep: init-dep-ui init-dep-api

.PHONY: init-dep-ui
init-dep-ui:
	@echo "> Initializing UI dependencies ..."
	@cd ${UI_PATH} && yarn

.PHONY: init-dep-api
init-dep-api:
	@echo "> Initializing API dependencies ..."
	@cd ${API_PATH} && go mod tidy -v
	@cd ${API_PATH} && go get -v ./...

# ============================================================
# Analyze source code recipes
# ============================================================
.PHONY: lint
lint: lint-ui lint-api

.PHONY: lint-ui
lint-ui:
	@echo "> Linting the UI source code ..."
	@cd ${UI_PATH} && yarn lint

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@goimports -w -local github.com/caraml-dev/mlp $(shell find . -type f -name '*.go' -not -path "**/vendor/*")
	@gofmt -s -w .

.PHONY: lint-api
lint-api: setup
	@echo "Linting code..."
	golangci-lint -v run --timeout 3m $(if $(filter true,$(fix)),--fix,) ${API_PATH}/...

# ============================================================
# Testing recipes
# ============================================================
.PHONY: test
test: test-api

.PHONY: test-api
test-api: init-dep-api
	@echo "> API unit testing ..."
	@cd ${API_PATH} && go test -v -race -cover -coverprofile cover.out ${API_ALL_PACKAGES}
	@cd ${API_PATH} && go tool cover -func cover.out

.PHONY: it-test-api-local
it-test-api-local: local-db start-keto start-vault
	@echo "> API integration testing locally..."
	@cd ${API_PATH} && go test -race -short -cover -coverprofile cover.out -tags integration ${API_ALL_PACKAGES}
	@cd ${API_PATH} && go tool cover -func cover.out
	@make stop-docker

.PHONY: it-test-api-ci
it-test-api-ci:
	@echo "> API integration testing ..."
	@cd ${API_PATH} && go test -race -short -cover -coverprofile cover.out -tags integration ${API_ALL_PACKAGES}
	@cd ${API_PATH} && go tool cover -func cover.out

# ============================================================
# Building recipes
# ============================================================
.PHONY: build
build: build-ui build-api

.PHONY: build-ui
build-ui: clean-ui
	@echo "> Building UI static build ..."
	@cd ${UI_PATH} && yarn lib build && yarn app build

.PHONY: build-api
build-api: clean-bin
	@echo "> Building API binary ..."
	@cd ${API_PATH} && go build -o ../bin/${BIN_NAME} ./cmd/main.go

.PHONY: build-api-image
build-api-image: version
	@$(eval IMAGE_TAG = $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/,)${BIN_NAME}-api:${VERSION})
	@echo "Building docker image: ${IMAGE_TAG}"
	docker build . \
		--tag ${IMAGE_TAG} \
		--file api.Dockerfile

.PHONY: build-image
build-image: version
	@$(eval IMAGE_TAG = $(if $(DOCKER_REGISTRY),$(DOCKER_REGISTRY)/,)${BIN_NAME}:${VERSION})
	@echo "Building docker image: ${IMAGE_TAG}"
	docker build . \
		--build-arg MLP_API_IMAGE \
		--tag ${IMAGE_TAG} \
		--file full.Dockerfile

# ============================================================
# Run recipes
# ============================================================
.PHONY: run
run: build-api local-db
	@echo "> Running application ..."
	@./bin/${BIN_NAME} --config config-dev.yaml

.PHONY: start-ui
start-ui:
	@echo "> Starting UI app ..."
	@cd ${UI_PATH} && yarn start

# ============================================================
# Utility recipes
# ============================================================
.PHONY: clean
clean: clean-ui clean-bin

.PHONY: clean-ui
clean-ui:
	@echo "> Cleaning up existing UI static build ..."
	@test ! -e ${UI_BUILD_PATH} || rm -r ${UI_BUILD_PATH}

.PHONY: clean-bin
clean-bin:
	@echo "> Cleaning up existing compiled binary ..."
	@test ! -e bin || rm -r bin

generate-client:
	@echo "> Generating API client ..."
	@swagger-codegen generate -i static/swagger.yaml -l go -o client -DpackageName=client
	@goimports -l -w client

.PHONY: local-env
local-env: local-db start-keto start-vault

.PHONY: local-db
local-db:
	@echo "> Starting up DB ..."
	@docker-compose up -d postgres

.PHONY: start-keto
start-keto:
	@echo "> Starting up keto server ..."
	@docker-compose up -d keto

.PHONY: start-vault
start-vault:
	@echo "> Starting up vault server ..."
	@docker-compose up -d vault

.PHONY: stop-docker
stop-docker:
	@echo "> Stopping Docker compose ..."
	@docker-compose down

.PHONY: swagger-ui
swagger-ui:
	@echo "> Starting up Swagger UI ..."
	@docker-compose up -d swagger-ui

.PHONY: version
version:
	$(eval VERSION=$(if $(OVERWRITE_VERSION),$(OVERWRITE_VERSION),v$(shell ./scripts/vertagen/vertagen.sh)))
	@echo "version:" $(VERSION)
