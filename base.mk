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

GO_VET_TAGS ?= integration,containers_image_openpgp
vet: ##@Development Run go vet against code.
	go vet -tags $(GO_VET_TAGS) ./...

# you can set lint configuration by GOLANGCILINT_CONFIG like GOLANGCILINT_CONFIG=--issues-exit-code=0
GOLANGCILINT_CONFIG ?=
lint: golangcilint ##@Development Run golangci-lint against code.
	$(GOLANGCILINT) run $(GOLANGCILINT_CONFIG)

ENVTEST_ASSETS_DIR=$(TOOLBIN)/testbin
COVER_PROFILE ?= cover.out
TEST_FILE ?= test.json
GO_TEST_FLAGS ?= -v -json
GO_TEST_TAGS ?= containers_image_openpgp
gotest: ##@Development Run go tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); KUBEBUILDER_ASSETS=$(TOOLBIN)/testbin/bin go test $(GO_TEST_FLAGS) -tags=$(GO_TEST_TAGS) -coverpkg=./... -coverprofile $(COVER_PROFILE) ./... | tee ${TEST_FILE}

test: manifests generate fmt vet goimports gotest ##@Development Run source code tests.

install: manifests kustomize ##@Deployment Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ##@Deployment Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ko certmanager ##@Deployment Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) apply -P ${LOCAL} -f -

wait: manifests kustomize ko yq ##@Deployment Wait for deployment to complete
	$(KUSTOMIZE) build config/default | $(YQ) 'select(.kind == "Deployment") | .metadata | {.namespace:.name}' | grep -v -- --- | awk -F ': ' '{print "kubectl -n "$$1" rollout status deploy "$$2}' | sh

deploy-wait: deploy wait ##@Deployment Deploy controller to the K8s cluster and wait for completion

undeploy: kustomize ko ##@Deployment Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KO) delete -f -

certmanager: ##@Deployment Install certmanager v1.4.0 from github manifest to the K8s cluster specified in ~/.kube/config.
	$(call installyaml,cert-manager,https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml,cert-manager)

INTEGRATION_OPTIONS ?=
integration: ginkgo ##@Testing Executes integration tests inside test/integration folder
	$(GINKGO) -progress -v -tags $(GO_VET_TAGS) $(INTEGRATION_OPTIONS) ./test/integration

VULNCHECK_DB ?= https://vuln.go.dev
VULNCHECK_MODE ?= source
VULNCHECK_PATH ?= ./...
VULNCHECK_OUTPUT ?= vulncheck.txt
vulncheck: govulncheck ##@Development Run govulncheck against code. Check base.mk file for available envvars
	$(GOVULNCHECK) -db=$(VULNCHECK_DB) -mode=$(VULNCHECK_MODE) -tags $(GO_VET_TAGS) $(VULNCHECK_PATH) | tee $(VULNCHECK_OUTPUT)

TRIVY_DB_REPO ?= ghcr.io/aquasecurity/trivy-db
TRIVY_CACHE ?= $(HOME)/.cache/trivy
TRIVY_FORMAT ?= table
TRIVY_REPORT_OUTPUT ?= trivy-report.json
TRIVY_SCANNERS ?= vuln
TRIVY_SEVERITY ?= UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL
TRIVY_EXTRA_OPTIONS ?= --ignore-unfixed
trivy-repo-scan: trivy ##@Development Run trivy against code. Check base.mk file for available envvars
	$(TRIVY) repo  --vuln-type library \
		--db-repository=$(TRIVY_DB_REPO) \
		--cache-dir=$(TRIVY_CACHE) \
		--format=json \
		--output=$(TRIVY_REPORT_OUTPUT) \
		--scanners=$(TRIVY_SCANNERS) \
		--exit-code=0 \
		$(TRIVY_EXTRA_OPTIONS) .
	$(TRIVY) convert --format=$(TRIVY_FORMAT) \
		--severity=$(TRIVY_SEVERITY)  \
		--exit-code=1 $(TRIVY_REPORT_OUTPUT)

CRD_DOCS_SOURCE_PATH ?= pkg/apis
CRD_DOCS_OUTPUT_PATH ?= crd-docs
CRD_DOCS_MAX_DEPTH ?= 15
CRD_DOCS_CONFIG ?= .crd-docs.yaml
crd-docs: crd-ref-docs
	mkdir -p $(CRD_DOCS_OUTPUT_PATH)
	$(CRDREFDOCS) --config=$(CRD_DOCS_CONFIG) --source-path=$(CRD_DOCS_SOURCE_PATH) \
		--output-path=$(CRD_DOCS_OUTPUT_PATH) \
		--renderer=markdown \
		--output-mode=group \
		--max-depth=$(CRD_DOCS_MAX_DEPTH)

CONTROLLER_TOOLS_VERSION ?= v0.14.0
CONTROLLER_GEN = $(TOOLBIN)/controller-gen-$(CONTROLLER_TOOLS_VERSION)
controller-gen: ##@Setup Download controller-gen locally if necessary.
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

KUSTOMIZE_VERSION ?= v5.3.0
KUSTOMIZE = $(TOOLBIN)/kustomize-$(KUSTOMIZE_VERSION)
kustomize: ##@Setup Download kustomize locally if necessary.
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

KO_VERSION ?= v0.15.2
KO = $(TOOLBIN)/ko-$(KO_VERSION)
ko: ##@Setup Download ko locally if necessary.
	$(call go-install-tool,$(KO),github.com/google/ko,$(KO_VERSION))

GOIMPORTS_VERSION ?= v0.20.0
GOIMPORTS = $(TOOLBIN)/goimports-$(GOIMPORTS_VERSION)
goimports: ##@Setup Download goimports locally if necessary.
	$(call go-install-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))
	$(GOIMPORTS) -w -l $(shell find . -path '.git' -prune -path './vendor' -prune -o -path './examples' -prune -o -name '*.pb.go' -prune -o -type f -name '*.go' -print)

GINKGO_VERSION ?= v2.17.1
GINKGO = $(TOOLBIN)/ginkgo-$(GINKGO_VERSION)
ginkgo: ##@Setup Download ginkgo locally if necessary
	$(call go-install-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo,$(GINKGO_VERSION))

GOLANGCILINT_VERSION ?= v1.56.2
GOLANGCILINT = $(TOOLBIN)/golangci-lint-$(GOLANGCILINT_VERSION)
golangcilint: ##@Setup Download golangci-lint locally if necessary
	$(call go-install-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCILINT_VERSION))

YQ_VERSION ?= v4.43.1
YQ = $(TOOLBIN)/yq-$(YQ_VERSION)
yq: ##@Setup Download yq locally if necessary.
	$(call go-install-tool,$(YQ),github.com/mikefarah/yq/v4,$(YQ_VERSION))

GOMOCK_VERSION ?= v1.6.0
GOMOCK = $(TOOLBIN)/mockgen-$(GOMOCK_VERSION)
GOMOCK_ALIAS = $(TOOLBIN)/mockgen
gomock: ##@Setup Download gomock locally if necessary.
	$(call go-install-tool,$(GOMOCK),github.com/golang/mock/mockgen,$(GOMOCK_VERSION),$(GOMOCK_ALIAS))
	$(shell ln -s $(GOMOCK) $(GOMOCK_ALIAS))

APISERVER_RUNTIME_GEN_VERSION ?= v1.1.1
APISERVER_RUNTIME_GEN = $(TOOLBIN)/apiserver-runtime-gen-$(APISERVER_RUNTIME_GEN_VERSION)
apiserver-runtime-gen: ##@Setup Download apiserver-runtime-gen locally if necessary
	$(call go-install-tool,$(APISERVER_RUNTIME_GEN),sigs.k8s.io/apiserver-runtime/tools/apiserver-runtime-gen,$(APISERVER_RUNTIME_GEN_VERSION))

GOVULNCHECK_VERSION ?= master
GOVULNCHECK = $(TOOLBIN)/govulncheck-$(GOVULNCHECK_VERSION)
govulncheck: ##@Setup Download govulncheck locally if necessary.
# using master until 1.0.5 is released, https://github.com/golang/go/issues/66139
	$(call go-install-tool,$(GOVULNCHECK),golang.org/x/vuln/cmd/govulncheck,$(GOVULNCHECK_VERSION))

TRIVY_VERSION ?= 0.50.1
TRIVY = $(TOOLBIN)/trivy-$(TRIVY_VERSION)
trivy: ##@Setup Download trivy locally if necessary.
	$(call download-trivy,$(TRIVY),$(TRIVY_VERSION))

CRDREFDOCS_VERSION ?= v0.0.12
CRDREFDOCS = $(TOOLBIN)/crd-ref-docs
crd-ref-docs: ##@Setup Download crd-ref-docs locally if necessary.
	$(call go-install-tool,$(CRDREFDOCS),github.com/elastic/crd-ref-docs,$(CRDREFDOCS_VERSION))

githook: precommit ##@Development Install git pre-commit hook
	pre-commit install

precommit: ##@Setup Download pre-commit locally if necessary.
	pip install pre-commit

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package} into $(TOOLBIN) as $(1)" ;\
GOBIN=$(TOOLBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
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

# finds and downloads the binary for trivy from github releases
# $1 - target path to binary
# $2 - version without v i.e 0.50.1
define download-trivy
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
case $(shell uname -s) in \
  Darwin) 						OS=macOS ;; \
  Linux) 						OS=linux ;; \
  *) echo "Unsupported OS" >&2; exit 1 ;; \
esac ;\
case $(shell uname -m) in \
  arm64) ARCH=ARM64 ;; \
  x86_64) ARCH=64bit ;; \
  *) echo "Unsupported ARCH" >&2; exit 1 ;; \
esac ;\
echo "Downloading trivy $(2)" ;\
curl -L https://github.com/aquasecurity/trivy/releases/download/v$(2)/trivy_$(2)_$${OS}-$${ARCH}.tar.gz > download.tgz ;\
tar xzf download.tgz ;\
mv trivy $(1) ;\
rm -rf $$TMP_DIR ;\
}
endef
