# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

LOCAL ?=

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: test

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

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

lint: golangcilint ## Run golangci-lint against code.
	$(GOLANGCILINT) run

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet goimports ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

##@ Setup

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	## this is a necessary evil already reported by knative community https://github.com/kubernetes-sigs/controller-tools/ issue 560
	## once the issue is fixed we can move to use the original package. the original line uses go-get-tools with sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
	$(call go-get-fork,$(CONTROLLER_GEN),https://github.com/danielfbm/controller-tools,cmd/controller-gen,controller-gen)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

KO = $(shell pwd)/bin/ko
ko: ## Download ko locally if necessary.
	$(call go-get-tool,$(KO),github.com/google/ko@v0.8.3)

GOIMPORTS = $(shell pwd)/bin/goimports
goimports: ## Download goimports locally if necessary.
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)
	$(GOIMPORTS) -w -l $(shell find . -path '.git' -prune -path './vendor' -prune -o -path './examples' -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)

GINKGO = $(shell pwd)/bin/ginkgo
ginkgo: ## Download ginkgo locally if necessary
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo@v1.16.4)

GOMOCK = $(shell pwd)/bin/mockgen
gomock: ## Download gomock locally if necessary.
	$(call go-get-tool,$(GOMOCK),github.com/golang/mock/mockgen@v1.6.0)

GOLANGCILINT = $(shell pwd)/bin/golangci-lint
golangcilint: ## Download golangci-lint locally if necessary
	$(call go-get-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# go-get-fork is a "go-get-tool" like command to get temporary module forks.
define go-get-fork
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
echo "Cloning $(2)" ;\
git clone $(2) $(4) ;\
cd $(4) ;\
GOBIN=$(PROJECT_DIR)/bin go install ./$(3);\
rm -rf $$TMP_DIR ;\
}
endef

# installyaml will check if a given namespace is present, if not will apply a yaml file and wait for a deployment to rollout
define installyaml
kubectl get ns $(1) > /dev/null ;\
EXIT_CODE=$$?;\
[ "$$EXIT_CODE" == "0" ] || { \
set -e ;\
kubectl apply -f $(2) ;\
kubectl -n $(1) rollout status deploy/$(3) --timeout=10m ;\
}
endef
