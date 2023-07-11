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
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type CloudEvent struct {
	ID      string `json:"id,omitempty" variable:"example=b51b6a90be6a6f7a2aa65049ad.2022-08-05-05-34"`
	Source  string `json:"source,omitempty" variable:"example=https://github.com/repository"`
	Subject string `json:"subject,omitempty" variable:"example=58127"`
	// Type of event
	Type string `json:"type,omitempty" variable:"example=dev.katanomi.cloudevents.gitlab.Merge Request Hook"`
	// Data event payload
	Data            string            `json:"data,omitempty" variable:"-"`
	Time            metav1.Time       `json:"time,omitempty" variable:"example=2022-08-05T05:34:39Z"`
	SpecVersion     string            `json:"specversion,omitempty" variable:"example=1.0"`
	DataContentType string            `json:"datacontenttype,omitempty" variable:"example=application/json"`
	Extensions      map[string]string `json:"extensions,omitempty"`
}

func (evt *CloudEvent) From(event cloudevents.Event) *CloudEvent {
	evt.ID = event.ID()
	evt.Source = event.Source()
	evt.Data = string(event.Data())
	evt.Subject = event.Subject()
	evt.DataContentType = event.DataContentType()
	evt.Type = event.Type()
	evt.SpecVersion = event.SpecVersion()
	evt.Time = metav1.NewTime(event.Time())
	for key, val := range event.Extensions() {
		if evt.Extensions == nil {
			evt.Extensions = map[string]string{}
		}

		var str string

		switch v := val.(type) {
		case string:
			str = v
		case int, int8, int16, int32, int64:
			str = fmt.Sprintf("%d", v)
		case time.Time:
			str = v.String()
		case float32, float64:
			str = fmt.Sprintf("%d", v)
		case bool:
			str = fmt.Sprintf("%t", v)
		default:
			bts, _ := json.Marshal(v)
			str = string(bts)
		}
		evt.Extensions[key] = str
	}
	return evt
}

// GetValWithKey returns the list of keys and values to support variable substitution
func (evt *CloudEvent) GetValWithKey(ctx context.Context, path *field.Path) (values map[string]string) {
	if evt == nil {
		evt = &CloudEvent{}
	}
	values = map[string]string{
		path.String():                          "",
		path.Child("id").String():              evt.ID,
		path.Child("source").String():          evt.Source,
		path.Child("subject").String():         evt.Subject,
		path.Child("type").String():            evt.Type,
		path.Child("time").String():            evt.Time.UTC().Format(time.RFC3339),
		path.Child("specversion").String():     evt.SpecVersion,
		path.Child("datacontenttype").String(): evt.DataContentType,
		path.Child("extensions").String():      "",
	}
	for k, v := range evt.Extensions {
		values[path.Child("extensions", k).String()] = v
	}
	return
}

const (
	CloudEventPrefix                 = "dev.katanomi.cloudevents"
	CloudEventExtGitReference        = "reference"
	CloudEventExtGitBranch           = "branch"
	CloudEventExtGitCommitMessage    = "commitmessage"
	CloudEventExtGitPullRequestTitle = "pullrequesttitle"
	CloudEventExtGitCommitID         = "commit"
	CloudEventExtGitSourceBranch     = "sourcebranch"
	CloudEventExtGitTargetBranch     = "targetbranch"
	CloudEventExtGitTag              = "tag"
	CloudEventExtAction              = "action"
	CloudEventExtSender              = "sender"
	CloudEventExtPullRequestNumber   = "number"
	CloudEventExtCodeRepository      = "repository"
	CloudEventExtWebhookType         = "webhooktype"
	// CloudEventExtRevisionSubmitter indicates email of revision submitter
	CloudEventExtRevisionSubmitter = "revisionsubmitter"

	// Triggered when a pull request's head branch is updated.
	// For example, when the head branch is updated from the base branch, when new commits are pushed to the head branch, or when the base branch is changed.
	// Used for matching events in trigger filter
	CloudEventExtPullRequestActionSynchronize = "synchronize"
	// pull request title changed or description changed.
	// Used for matching events in trigger filter
	CloudEventExtPullRequestActionEdited = "edited"
	// Triggered when a pull request is closed
	// Used for matching events in trigger filter
	CloudEventExtPullRequestActionClosed = "closed"
	// Triggered when a pull request is merged
	// Used for matching events in trigger filter
	CloudEventExtPullRequestActionMerged = "merged"
	// event type of open pull request
	CloudEventExtPullRequestActionOpened = "opened"
	// event type of reopen pull request
	CloudEventExtPullRequestActionReOpened = "reopened"

	// action of create branch
	// Used for matching events in trigger filter
	CloudEventExtBranchActionCreate = "create"
	// action of delete branch
	// Used for matching events in trigger filter
	CloudEventExtBranchActionDelete = "delete"

	// action of create tag
	// Used for matching events in trigger filter
	CloudEventExtTagActionCreate = "create"
	// action of delete tag
	// Used for matching events in trigger filter
	CloudEventExtTagActionDelete = "delete"

	CloudEventExtArtifactActionPush = "push"
	// artifact type. eg. OCIHelmChart OCIContainerImage
	CloudEventExtArtifactType = "artifacttype"
	// resource URL. eg. katanomi/core:latest
	CloudEventExtArtifactResourceURL = "resourceurl"
	// is the digest of the manifest. eg. sha256:92d648...
	CloudEventExtArtifactDigest = "digest"
	// is the tag of the artifact. (if not assigned, the value same as digest)
	CloudEventExtArtifactTag = "tag"
	// username of user who pushed artifact
	CloudEventExtArtifactOperator = "operator"
	// push time timestamp
	CloudEventExtArtifactOccurAt = "occurat"
	// concat os-arch-variant and os-variant strings of all artifact, use # connect each string.
	// eg. #linux-amd64-#linux-#linux-arm64-v8#linux-v8# (only supported when artifact_type is OCIContainerImage)
	CloudEventExtArtifactPlatform = "platform"
)
