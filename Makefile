# Image URL to use all building/pushing image targets
IMG ?= kubectl-onmetal:latest

.PHONY: all
all: check

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: goimports ## Run goimports against code.
	$(GOIMPORTS) -w .

.PHONY: add-license
add-license: addlicense ## Add license headers to all go files.
	find . -name '*.go' -exec $(ADDLICENSE) -c 'OnMetal authors' {} +

.PHONY: check-license
check-license: addlicense ## Check that every file has a license header present.
	find . -name '*.go' -exec $(ADDLICENSE) -check -c 'OnMetal authors' {} +

.PHONY: generate
generate: ## Generate code (mocks etc.).
	go generate ./...

.PHONY: test
test: generate fmt ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: lint
lint: ## Run the linter.
	golangci-lint run ./...

.PHONY: check
check: generate add-license fmt lint test ## Execute multiple checks.

.PHONY: build
build: generate fmt ## Build the binary
	go build -o bin/kubectl-onmetal .

.PHONY: install
install: generate fmt ## Install the binary to /usr/local/bin
	go install .

.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Tools

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ADDLICENSE ?= $(LOCALBIN)/addlicense
GOIMPORTS ?= $(LOCALBIN)/goimports

## Tool Versions
ADDLICENSE_VERSION ?= v1.1.0
GOIMPORTS_VERSION ?= v0.5.0

.PHONY: addlicense
addlicense: $(ADDLICENSE) ## Download addlicense locally if necessary.
$(ADDLICENSE): $(LOCALBIN)
	test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: goimports
goimports: $(GOIMPORTS) ## Download goimports locally if necessary.
$(GOIMPORTS): $(LOCALBIN)
	test -s $(LOCALBIN)/goimports || GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

