/*
Copyright 2021 The Katanomi Authors.

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

package v1alpha1

import (
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

const (
	// indicates trigger name on subscription
	TriggerLabel string = "core.katanomi.dev/trigger"
	// indicates git trigger name on subscription
	GitTriggerLabel string = "core.katanomi.dev/gitTrigger"
	// indicates artifact trigger name on subscription
	ArtifactTriggerLabel string = "core.katanomi.dev/artifactTrigger"
	// indicates broker name on subscription
	BrokerLabel string = "eventing.knative.dev/broker"
	// indicates integration class name
	IntegrationClassLabel = metav1alpha1.IntegrationClassLabelKey
)

const (
	// ReconcileTriggeredAnnotation adds an annotations to trigger reconcile of an object
	ReconcileTriggeredAnnotation = metav1alpha1.ReconcileTriggeredAnnotationKey
	// WebhookAutoGenerateAnnotation indicates that will auto generate webhook
	WebhookAutoGenerateAnnotation string = "core.katanomi.dev/webhook.autoGenerate"
	// WebhookOwnerHashAnnotation indicates the webhook owner's hash code, if changed, should be reconciled.
	WebhookOwnerHashAnnotation string = "core.katanomi.dev/webhook.owner.hash"
)
