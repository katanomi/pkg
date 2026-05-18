# Pre-include overrides: temporary forward-port of fixes that have not yet
# auto-synced into this repo's base.mk from katanomi/hack. Each uses `?=`
# (not `=`) so env/CI overrides still win, matching base.mk's idiom. Once
# base.mk catches up these overrides become no-ops.
#
# TRIVY_VERSION: base.mk pins 0.50.1, but aquasecurity/trivy has no v0.50.1
# release (timeline jumps from v0.26 to v0.69.x); CI's trivy-repo-scan task
# 404s on download. Tracks katanomi/hack#66.
TRIVY_VERSION ?= 0.69.3
#
# CONTROLLER_TOOLS_VERSION: base.mk pins v0.14.0, which pulls
# golang.org/x/tools v0.16.1 transitively. That x/tools fails to compile
# under Go 1.25.x (`tokeninternal.go:78: invalid array length -delta * delta`),
# breaking `make controller-gen` and therefore `make test`. v0.17.2 matches
# the prebuilt binary shipped by the rebuilt builder-go:latest image.
CONTROLLER_TOOLS_VERSION ?= v0.17.2
#
# GOLANGCILINT_VERSION: base.mk pins v1.56.2; under Go 1.25.x with the new
# dependency graph the linter both warns "version too low, upgrade to 1.57+"
# and times out with "context deadline exceeded". v1.62.2 is the first 1.6x
# series that handles Go 1.25 cleanly in our CI.
GOLANGCILINT_VERSION ?= v1.62.2

include base.mk

manifests: controller-gen ##@Development Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	# $(CONTROLLER_GEN) rbac:roleName=pkg paths="./..."
