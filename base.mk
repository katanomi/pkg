###################################################
### WARNING: This file is synced from katanomi/hack
### DO NOT CHANGE IT MANUALLY
###################################################
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd"

LOCAL ?=

TOOLBIN ?= $(shell pwd)/bin

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

all: check

HELP_FUN = \
	%help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
	if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
	print"\033[1m$$_:\033[0m\n", map"  \033[36m$$_->[0]\033[0m".(" "x(20-length($$_->[0])))."$$_->[1]\n",\
	@{$$help{$$_}},"\n" for keys %help; \

help: ##@General Show this help
	@echo -e "Usage: make \033[36m<target>\033[0m\n"
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

generate: controller-gen ##@Development Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

check: fmt vet lint test ##@Development Run check against code

fmt: ##@Development Run go fmt against code.
	go fmt ./...

GO_VET_TAGS ?= e2e,containers_image_openpgp
vet: ##@Development Run go vet against code.
	go vet -tags $(GO_VET_TAGS) ./...

lint: golangcilint ##@Development Run golangci-lint against code.
	$(GOLANGCILINT) run

ENVTEST_ASSETS_DIR=$(TOOLBIN)/testbin
COVER_PROFILE ?= cover.out
TEST_FILE ?= test.json
GO_TEST_FLAGS ?= -v -json
GO_TEST_TAGS ?= containers_image_openpgp
gotest: ##@Development Run go tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test $(GO_TEST_FLAGS) -tags=$(GO_TEST_TAGS) -coverpkg=./... -coverprofile $(COVER_PROFILE) ./... | tee ${TEST_FILE}

test: manifests generate fmt vet goimports gotest ##@Development Run source code tests.

install: manifests kustomize ##@Deployment Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ##@Deployment Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ko certmanager ##@Deployment Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) apply -P ${LOCAL} -f -

wait: manifests kustomize yq ##@Deployment Wait for deployment to complete
	$(KUSTOMIZE) build config/default | $(YQ) 'select(.kind == "Deployment") | .metadata | {.namespace:.name}' | grep -v -- --- | awk -F ': ' '{print "kubectl -n "$$1" rollout status deploy "$$2}' | sh

deploy-wait: deploy wait ##@Deployment Deploy controller to the K8s cluster and wait for completion

undeploy: ##@Deployment Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) delete -f -

certmanager: ##@Deployment Install certmanager v1.4.0 from github manifest to the K8s cluster specified in ~/.kube/config.
	$(call installyaml,cert-manager,https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml,cert-manager)

e2e: ginkgo ##@Testing Executes e2e tests inside test/e2e folder
	$(GINKGO) -progress -v -tags $(GO_VET_TAGS) ./test/e2e

CONTROLLER_GEN = $(TOOLBIN)/controller-gen
controller-gen: ##@Setup Download controller-gen locally if necessary.
	## this is a necessary evil already reported by knative community https://github.com/kubernetes-sigs/controller-tools/ issue 560
	## once the issue is fixed we can move to use the original package. the original line uses go-get-tools with sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
	$(call go-get-fork,$(CONTROLLER_GEN),https://github.com/danielfbm/controller-tools,cmd/controller-gen,controller-gen)

KUSTOMIZE = $(TOOLBIN)/kustomize
kustomize: ##@Setup Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4@v4.5.7)

KO = $(TOOLBIN)/ko
ko: ##@Setup Download ko locally if necessary.
	$(call go-get-tool,$(KO),github.com/google/ko@v0.12.0)

GOIMPORTS = $(TOOLBIN)/goimports
goimports: ##@Setup Download goimports locally if necessary.
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports@v0.1.10)
	$(GOIMPORTS) -w -l $(shell find . -path '.git' -prune -path './vendor' -prune -o -path './examples' -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)

GINKGO = $(TOOLBIN)/ginkgo
ginkgo: ##@Setup Download ginkgo locally if necessary
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo@v1.16.4)

GOLANGCILINT = $(TOOLBIN)/golangci-lint
golangcilint: ##@Setup Download golangci-lint locally if necessary
	$(call go-get-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2)

YQ = $(TOOLBIN)/yq
yq: ##@Setup Download yq locally if necessary.
	$(call go-get-tool,$(YQ),github.com/mikefarah/yq/v4@v4.25.2)

GOMOCK = $(TOOLBIN)/mockgen
gomock: ## Download gomock locally if necessary.
	$(call go-get-tool,$(GOMOCK),github.com/golang/mock/mockgen@v1.6.0)

apiserver-runtime-gen:
	$(call go-get-tool,$(GOMOCK),sigs.k8s.io/apiserver-runtime/tools/apiserver-runtime-gen@v1.1.1)

githook: precommit ##@Development Install git pre-commit hook
	pre-commit install

precommit: ##@Setup Download pre-commit locally if necessary.
	pip install pre-commit

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
GOBIN=$(TOOLBIN) go install $(2) ;\
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
GOBIN=$(TOOLBIN) go install ./$(3);\
rm -rf $$TMP_DIR ;\
}
endef

# installyaml will check if a given namespace is present, if not will apply a yaml file and wait for a deployment to rollout
define installyaml
kubectl get ns $(1) > /dev/null || { \
set -e ;\
kubectl apply -f $(2) ;\
kubectl -n $(1) rollout status deploy/$(3) --timeout=10m ;\
}
endef

