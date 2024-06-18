/*
Copyright 2023 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

const (
	// VersionEnabledFeatureKey indicates the configuration key of the version feature gate.
	// If the value is true, the feature is enabled cluster-wide.
	VersionEnabledFeatureKey = "version.enabled"

	// ProxyEnabledFeatureKey indicates the configuration key of the proxy feature gate.
	// If the value is true, the feature is enabled cluster-wide.
	ProxyEnabledFeatureKey = "proxy.enabled"

	// InitializeAllowLocalRequestsFeatureKey indicates the configuration key of.
	// If the value is true, the feature is enabled cluster-wide.
	InitializeAllowLocalRequestsFeatureKey = "plugin.gitlab.allow-local-requests"

	// PrunerDelayAfterCompletedFeatureKey represent taskRun delay configuration key
	PrunerDelayAfterCompletedFeatureKey = "taskRunPruner.delayAfterCompleted"

	// PrunerKeepFeatureKey represent taskRun keep configuration key
	PrunerKeepFeatureKey = "taskRunPruner.keep"

	// BuildMRCheckTimeoutKey represent build merge request status check timeout key
	BuildMRCheckTimeoutKey = "build.mergerequest.checkTimeout"

	// BuildMRCheckTimeoutContinueKey represent build merge request status check timeout continue key
	BuildMRCheckTimeoutContinueKey = "build.mergerequest.checkTimeoutContinue"

	// TemplateRenderCheckTimeoutKey represent templatrender timeout key
	TemplateRenderCheckTimeoutKey = "templateRender.checkTimeout"

	// TemplateRenderRetentionTimeKey represent the config key for how long templatrender will remain
	// retentionTime seems like a better name than delayAfterCompleted so that we use rententionTime here
	TemplateRenderRetentionTimeKey = "templateRender.retentionTime"

	// PolicyRunRetentionTimeKey represent the config key for how long policyRun will remain
	PolicyRunRetentionTimeKey = "policyRun.retentionTime"

	// PolicyCheckEnabledFeatureKey indicates the configuration key of the policy check feature gate.
	// If the value is true, the feature is enabled cluster-wide.
	PolicyCheckEnabledFeatureKey = "policy.check.enabled"

	// ClusterIntegrationSyncPeriodKey indicates the configuration key for clusterintegration synchronization.
	ClusterIntegrationSyncPeriodKey = "clusterIntegration.syncPeriod"

	// IntegrationSyncPeriodKey indicates the configuration key for integration synchronization.
	IntegrationSyncPeriodKey = "integration.syncPeriod"

	// IntegrationResourcesSyncPeriodKey indicates the configuration key for integration resources synchronization.
	// if the values is zero , it indicates never requeue in period
	IntegrationResourcesSyncPeriodKey = "integration.resources.syncPeriod"

	// PprofEnabledKey indicates the configuration key of the /debug/pprof debugging api/
	PprofEnabledKey = "pprof.enabled"

	// ClusterTaskDisabledKey specifies the key for the clustertask creation feature.
	// When set to false, the creation of clustertasks is disabled and new clustertasks cannot be created and can only be updated.
	ClusterTaskCreationEnabledKey = "clustertask.creation.enabled"

	// ClusterTaskMigrationEnabledKey specifies the key for the clustertask migration feature.
	ClusterTaskMigrationEnabledKey = "clustertask.migration.enabled"

	// KatanomiSystemNamespaceKey specifies the key for the katanomi system namespace
	KatanomiSystemNamespaceKey = "katanomi.system.namespace"

	// GitSourceResolverFeatureKey indicates the configuration key of the gitsource resolver feature gate.
	GitSourceResolverFeatureKey = "gitsourceResolver.enabled"

	// HubResolverFeatureKey indicates the configuration key of the hub resolver feature gate.
	HubResolverFeatureKey = "hubResolver.enabled"

	// GlobalCredentialNamespaceKey is the key of the global credentials namespace
	GlobalCredentialNamespaceKey = "globalCredentials.namespace"
)

const (
	// True represents the value "true" for the feature switch.
	True FeatureValue = "true"

	// False represents the value "false" for the feature switch.
	False FeatureValue = "false"

	// DefaultVersionEnabled indicates the default value of the version feature gate.
	// If the corresponding key does not exist, the default value is returned.
	DefaultVersionEnabled FeatureValue = False

	// DefaultProxyEnabled indicates the default value of the proxy feature gate.
	// If the corresponding key does not exist, the default value is returned.
	DefaultProxyEnabled FeatureValue = False

	// DefaultInitializeAllowLocalRequests indicates the configuration key of.
	// If the corresponding key does not exist, the default value is returned.
	DefaultInitializeAllowLocalRequests FeatureValue = True

	// DefaultPrunerDelayAfterCompleted represent default duration for delay taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerDelayAfterCompleted FeatureValue = "time.Hour"

	// DefaultPrunerKeep represent default keep number for taskRun
	// If the corresponding key does not exist, the default value is returned.
	DefaultPrunerKeep FeatureValue = "10000"

	// DefaultMRCheckTimeout represent default timeout for merge request status check
	DefaultMRCheckTimeout FeatureValue = "10m"

	// DefaultMRCheckTimeoutContinue represent default timeout continue for merge request status check
	DefaultMRCheckTimeoutContinue FeatureValue = True

	// DefaultTemplateRenderCheckTimeout represent default timeout for templaterender check
	DefaultTemplateRenderCheckTimeout FeatureValue = "30s"

	// DefaultTemplateRenderRetentionTime represents default duration how long the templatrender will remain
	DefaultTemplateRenderRetentionTime FeatureValue = "30m"

	// DefaultPolicyRunRetentionTime represents default duration how long the policyRun will remain
	DefaultPolicyRunRetentionTime = "30m"

	// DefaultPolicyCheckEnabled indicates the default value of the policy check feature gate.
	// If the corresponding key does not exist, the default value is returned.
	DefaultPolicyCheckEnabled FeatureValue = True

	// DefaultClusterIntegrationSyncPeriod defines the default time interval of clusterintegration synchronization
	DefaultClusterIntegrationSyncPeriod = "5m"

	// DefaultIntegrationsSyncPeriod defines the default time interval of integration synchronization
	DefaultIntegrationsSyncPeriod = "15m"

	// DefaultIntegrationsResourcesSyncPeriod defines the default time interval of integration synchronization for resources
	DefaultIntegrationsResourcesSyncPeriod = "2h"

	// DefaultPprofEnabled stores the default value "false" for the "pprof.enabled" /debug/pprof debugging api.
	// If the corresponding key does not exist, the default value is returned.
	DefaultPprofEnabled FeatureValue = False

	// DefaultClusterTaskCreationEnabledKey stores the default value "false" for the feature.
	DefaultClusterTaskCreationEnabledKey FeatureValue = False

	// DefaultClusterTaskMigrationEnabledKey stores the default value "false" for the feature.
	DefaultClusterTaskMigrationEnabledKey FeatureValue = False

	// DefaultKatanomiSystemNamespace defines the default namespace for katanomi system
	DefaultKatanomiSystemNamespace = "katanomi-system"

	// DefaultGitSourceResolverEnabled defines the default value "true" for the "gitsourceResolver.enabled" feature.
	DefaultGitSourceResolverEnabled = True

	// DefaultHubResolverEnabled defines the default value "true" for the "hubResolver.enabled" feature.
	DefaultHubResolverEnabled = True
)

// defaultFeatureValue defines the default value for the feature switch.
var defaultFeatureValue = map[string]FeatureValue{
	VersionEnabledFeatureKey:               DefaultVersionEnabled,
	ProxyEnabledFeatureKey:                 DefaultProxyEnabled,
	InitializeAllowLocalRequestsFeatureKey: DefaultInitializeAllowLocalRequests,
	PrunerDelayAfterCompletedFeatureKey:    DefaultPrunerDelayAfterCompleted,
	PrunerKeepFeatureKey:                   DefaultPrunerKeep,
	BuildMRCheckTimeoutKey:                 DefaultMRCheckTimeout,
	BuildMRCheckTimeoutContinueKey:         DefaultMRCheckTimeoutContinue,
	TemplateRenderCheckTimeoutKey:          DefaultTemplateRenderCheckTimeout,
	TemplateRenderRetentionTimeKey:         DefaultTemplateRenderRetentionTime,
	PolicyRunRetentionTimeKey:              DefaultPolicyRunRetentionTime,
	PolicyCheckEnabledFeatureKey:           DefaultPolicyCheckEnabled,
	ClusterIntegrationSyncPeriodKey:        DefaultClusterIntegrationSyncPeriod,
	IntegrationSyncPeriodKey:               DefaultIntegrationsSyncPeriod,
	IntegrationResourcesSyncPeriodKey:      DefaultIntegrationsResourcesSyncPeriod,
	PprofEnabledKey:                        DefaultPprofEnabled,
	ClusterTaskCreationEnabledKey:          DefaultClusterTaskCreationEnabledKey,
	ClusterTaskMigrationEnabledKey:         DefaultClusterTaskMigrationEnabledKey,
	KatanomiSystemNamespaceKey:             DefaultKatanomiSystemNamespace,
	GitSourceResolverFeatureKey:            DefaultGitSourceResolverEnabled,
	HubResolverFeatureKey:                  DefaultHubResolverEnabled,
}

// FeatureFlags holds the features configurations
type FeatureFlags struct {
	Data map[string]string
}

// FeatureValue returns the value of the implemented feature flag, or the default if not found.
func (f *FeatureFlags) FeatureValue(flag string) FeatureValue {
	defaultValue := defaultFeatureValue[flag]
	if f == nil || f.Data == nil {
		return defaultValue
	}

	if value, ok := f.Data[flag]; ok {
		return FeatureValue(value)
	}
	return defaultValue
}
