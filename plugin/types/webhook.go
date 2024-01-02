/*
Copyright 2024 The Katanomi Authors.

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

package types

import (
	"context"

	cloudevent "github.com/cloudevents/sdk-go/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

// WebhookRegister used to register and manage webhooks
type WebhookRegister interface {
	// Use the methods below to manage webhooks in the target platform
	CreateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	UpdateWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) (metav1alpha1.WebhookRegisterStatus, error)
	DeleteWebhook(ctx context.Context, spec metav1alpha1.WebhookRegisterSpec, secret corev1.Secret) error
	ListWebhooks(ctx context.Context, uri apis.URL, secret corev1.Secret) ([]metav1alpha1.WebhookRegisterStatus, error)
}

// WebhookResourceDiffer used to compare different webhook resources in order to provide
// a way to merge webhook registration requests. If not provided, the resource's URI will be directly compared
type WebhookResourceDiffer interface {
	// IsSameResource will provide two ResourceURI
	// the plugin should discern if they are the same.
	// If this method is not implemented a standard comparisons will be used
	IsSameResource(ctx context.Context, i, j metav1alpha1.ResourceURI) bool
}

// WebhookReceiver receives a webhook request with validation and transform it into a cloud event
type WebhookReceiver interface {
	Interface
	ReceiveWebhook(ctx context.Context, req *restful.Request, secret string) (cloudevent.Event, error)
}

// GitTriggerRegister used to register GitTrigger
// TODO: need refactor: maybe integration plugin should decided how to generate cloudevents filters
// up to now, it is not a better solution that relying on plugins to give some events type to GitTriggerReconcile.
//
// PullRequestCloudEventFilter() CloudEventFilters
// BranchCloudEventFilter() CloudEventFilters
// TagCloudEventFilter() CloudEventFilters
// WebHook() WebHook
type GitTriggerRegister interface {
	GetIntegrationClassName() string

	// cloud event type of pull request hook that will match
	PullRequestEventType() string

	// cloud event type of push hook that will match
	PushEventType() string

	// cloud event type of push hook that will match
	TagEventType() string
}
