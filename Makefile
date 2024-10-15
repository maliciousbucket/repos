ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: test
test: ## Run unit tests
	go test ./...

##@ Build

.PHONY: build
build: ## Build binary
	go build -o bin/repos .

.PHONY: docker-build
docker-build: ## Build Dockerfile
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push image to a container registry
	docker push ${IMG}

.PHONY: docker-buildx
docker-buildx: ## Build and push the image to a container registry
	docker buildx --push --tag ${IMG} -f Dockerfile

.PHONY: run
run: fmt vet ## Run the application locally
	go run .

##@ Workflows

.PHONY: workflow
workflow: ## Trigger a workflow by passing in an action
	@bash tools/trigger-workflow.sh ${ACTION} ${ENV}


.PHONY: workflow-ci
workflow-ci: ## Trigger the CI workflow
	@bash tools/trigger-workflow.sh ci ${ENV}

.PHONY: workflow-deploy
workflow-deploy: ## Trigger the Deployment workflow
	@bash tools/trigger-workflow.sh deployment ${ENV}

.PHONY: workflow-test
workflow-test: ## Trigger the Test Run workflow
	@bash  tools/trigger-workflow.sh test-run ${ENV}






.PHONY: help
help:
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif