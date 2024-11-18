include base.mk

manifests: controller-gen ##@Development Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	# $(CONTROLLER_GEN) rbac:roleName=pkg paths="./..."

.PHONY: htmlreport
htmlreport: go-test-report##@Development open test report
	cat test.json | $(GO_TEST_HTML_REPORT)

GO_TEST_HTML_REPORT ?= $(TOOLBIN)/go-test-report
GO_TEST_HTML_REPORT_VERSION ?= v0.9.3
go-test-report: ##@Development install go-test-report
	$(call go-install-tool,$(GO_TEST_HTML_REPORT),github.com/vakenbolt/go-test-report,$(GO_TEST_HTML_REPORT_VERSION))
