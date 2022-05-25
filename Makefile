include base.mk

manifests: controller-gen ##@Development Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	# $(CONTROLLER_GEN) rbac:roleName=pkg paths="./..."
