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
	"io"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// ArtifactLister list artifact
type ArtifactLister interface {
	Interface
	ListArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error)
}

// ArtifactGetter get artifact detail
type ArtifactGetter interface {
	Interface
	GetArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) (*metav1alpha1.Artifact, error)
}

// ArtifactDeleter delete artifact
type ArtifactDeleter interface {
	Interface
	DeleteArtifact(ctx context.Context, params metav1alpha1.ArtifactOptions) error
}

// ProjectArtifactLister list project-level artifacts
type ProjectArtifactLister interface {
	Interface
	ListProjectArtifacts(ctx context.Context, params metav1alpha1.ArtifactOptions, option metav1alpha1.ListOptions) (*metav1alpha1.ArtifactList, error)
}

// ProjectArtifactGetter get artifact detail
type ProjectArtifactGetter interface {
	Interface
	GetProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (*metav1alpha1.Artifact, error)
}

// ProjectArtifactDeleter delete artifact
type ProjectArtifactDeleter interface {
	Interface
	DeleteProjectArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) error
}

// ProjectArtifactUploader upload artifact
type ProjectArtifactUploader interface {
	Interface
	UploadArtifact(ctx context.Context, params metav1alpha1.ProjectArtifactOptions, r io.Reader) error
}

// ProjectArtifactFileGetter download artifact within a project
type ProjectArtifactFileGetter interface {
	Interface
	GetProjectArtifactFile(ctx context.Context, params metav1alpha1.ProjectArtifactOptions) (io.ReadCloser, error)
}

// ArtifactTagDeleter delete a specific tag of the artifact.
type ArtifactTagDeleter interface {
	Interface
	DeleteArtifactTag(ctx context.Context, params metav1alpha1.ArtifactTagOptions) error
}

// ArtifactTriggerRegister used to register ArtifactTrigger
type ArtifactTriggerRegister interface {
	GetIntegrationClassName() string

	// cloud event type of push hook that will match
	PushEventType() string
}
