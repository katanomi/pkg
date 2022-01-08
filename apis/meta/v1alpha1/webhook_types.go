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
	"knative.dev/pkg/apis"
)

// WebhookEventSupportType consists of event type and "Event"
type WebhookEventSupportType string

func (w WebhookEventSupportType) String() string {
	return string(w)
}

const (
	// https://github.com/katanomi/spec/pull/80
	// key in integrationClass.status
	// CodeRepositoryPushWebhookEvent event value for code repository's push webhook
	CodeRepositoryPushWebhookEvent WebhookEventSupportType = "CodeRepositoryPushEvent"
	// CodeRepositoryTagWebhookEvent event values for code repository's tag webhook
	CodeRepositoryTagWebhookEvent WebhookEventSupportType = "CodeRepositoryTagEvent"
	// CodeRepositoryPullRequestWebhookEvent event values for code repository's pr webhook
	CodeRepositoryPullRequestWebhookEvent WebhookEventSupportType = "CodeRepositoryPullRequestEvent"
	// ArtifactDeleteWebhookEvent event values for artifact's delete webhook
	ArtifactDeleteWebhookEvent WebhookEventSupportType = "ArtifactDeleteEvent"
	// ArtifactPushWebhookEvent event values for artifact's push webhook
	ArtifactPushWebhookEvent WebhookEventSupportType = "ArtifactPushEvent"

	// WebhookEventSuffix key's suffix in integrationClass.status
	WebhookEventSuffix = "Event"
)

// WebhookRegisterSpec specifications to register a webhook
type WebhookRegisterSpec struct {
	// URI stores the Uniform Resource Identifier for webhook resource
	URI apis.URL `json:"uri"`
	// Events holds a list of event types desired for registration
	Events []string `json:"events"`
	// WebhookID is only used in update and can be blank for creation
	// +optional
	WebhookID string `json:"webhookID"`
	// Address stores the target address for the webhook
	Address apis.URL `json:"addressURL"`
	// RequestSecret will hold information for a request header that should be used
	// by the registring webhook
	// this data will be used during request to validate webhook requests
	// +optional
	RequestSecret string `json:"requestSecret"`
}

// WebhookRegisterStatus stores a registration request result
type WebhookRegisterStatus struct {
	// WebhookID will return the registered webhook id during create
	WebhookID string `json:"webhookID"`

	// Address stores the target address for the webhook
	Address apis.URL `json:"addressURL"`

	// StatusCode for the API request
	// +optional
	StatusCode int `json:"statusCode"`

	// Body response body returned from the request
	// +optional
	Body []byte `json:"body"`
}
