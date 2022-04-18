# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd"

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

vet: ##@Development Run go vet against code.
	go vet ./...

lint: golangcilint ##@Development Run golangci-lint against code.
	$(GOLANGCILINT) run

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet goimports ##@Development Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

install: manifests kustomize ##@Deployment Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ##@Deployment Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ko certmanager ##@Deployment Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) apply -P ${LOCAL} -f -

undeploy: ##@Deployment Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) delete -f -

certmanager: ##@Deployment Install certmanager v1.4.0 from github manifest to the K8s cluster specified in ~/.kube/config.
	$(call installyaml,cert-manager,https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml,cert-manager)

e2e: ginkgo ##@Testing Executes e2e tests inside test/e2e folder
	$(GINKGO) -progress -v -tags e2e ./test/e2e

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ##@Setup Download controller-gen locally if necessary.
	## this is a necessary evil already reported by knative community https://github.com/kubernetes-sigs/controller-tools/ issue 560
	## once the issue is fixed we can move to use the original package. the original line uses go-get-tools with sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
	$(call go-get-fork,$(CONTROLLER_GEN),https://github.com/danielfbm/controller-tools,cmd/controller-gen,controller-gen)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ##@Setup Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

KO = $(shell pwd)/bin/ko
ko: ##@Setup Download ko locally if necessary.
	$(call go-get-tool,$(KO),github.com/google/ko@v0.8.3)

GOIMPORTS = $(shell pwd)/bin/goimports
goimports: ##@Setup Download goimports locally if necessary.
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)
	$(GOIMPORTS) -w -l $(shell find . -path '.git' -prune -path './vendor' -prune -o -path './examples' -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)

GINKGO = $(shell pwd)/bin/ginkgo
ginkgo: ##@Setup Download ginkgo locally if necessary
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo@v1.16.4)

GOLANGCILINT = $(shell pwd)/bin/golangci-lint
golangcilint: ##@Setup Download golangci-lint locally if necessary
	$(call go-get-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2)

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

