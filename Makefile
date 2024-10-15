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

.PHONY: fmt ## Run go fmt against code
fmt:
	go fmt ./...

.PHONY: vet ## Run go vet against code
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

##@ Build

.PHONY: build ## Build binary
build:
	go build -o bin/repos .

.PHONY: docker-build # Build Dockerfile
docker-build:
	docker build -t ${IMG} .

.PHONY: docker-push # Push image to a container registry
docker-push:
	docker push ${IMG}

.PHONY: docker-buildx # Build and push the image to a container registry
docker-buildx:
	docker buildx --push --tag ${IMG} -f Dockerfile

.PHONY: run # Run the application locally
run: fmt vet
	go run .

##@ Workflows

.PHONY: workflow # Trigger a workflow by passing in an action
workflow:
	$(shell) tools/trigger-workflows.sh ${ACTION} ${ENV}


.PHONY: workflow-ci # Trigger the CI workflow
workflow:
	$(shell) tools/trigger-workflows.sh ci ${ENV}

.PHONY: workflow-deployment # Trigger the deployment workflow
workflow:
	$(shell) tools/trigger-workflows.sh deployment ${ENV}

.PHONY: workflow-test # Trigger the Test Run workflow
workflow:
	$(shell) tools/trigger-workflows.sh test-run ${ENV}






.PHONY: help
help:
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif