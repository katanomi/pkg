
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

test: fmt vet gomock goimports ## Run go test ./... against code
	go test ./... -coverprofile cover.out

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	## this is a necessary evil already reported by knative community https://github.com/kubernetes-sigs/controller-tools/ issue 560
	## once the issue is fixed we can move to use the original package. the original line uses go-get-tools with sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
	$(call go-get-fork,$(CONTROLLER_GEN),https://github.com/danielfbm/controller-tools,cmd/controller-gen,controller-gen)

GOMOCK = $(shell pwd)/bin/gomock
gomock: ## Download gomock locally if necessary.
	$(call go-get-tool,$(GOMOCK),github.com/golang/mock/mockgen@v1.6.0)

GOIMPORTS = $(shell pwd)/bin/goimports
goimports: ## Download goimports locally if necessary.
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)
	$(GOIMPORTS) -w -v -l $(shell find . -path '.git' -prune -path './vendor' -prune -o -path './examples' -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)

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
