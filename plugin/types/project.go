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

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

//go:generate mockgen -package=types -destination=../../testing/mock/github.com/katanomi/pkg/plugin/types/project.go github.com/katanomi/pkg/plugin/types ProjectLister,ProjectGetter,SubtypeProjectGetter,ProjectCreator,ProjectDeleter

// ProjectLister list project api
type ProjectLister interface {
	Interface
	ListProjects(ctx context.Context, option metav1alpha1.ListOptions) (*metav1alpha1.ProjectList, error)
}

// ProjectGetter list project api
// Deprecated: use ProjectSubTypeGetter instead
type ProjectGetter interface {
	Interface
	GetProject(ctx context.Context, id string) (*metav1alpha1.Project, error)
}

// GetProjectOption option to get a subtype project
type GetProjectOption struct {
	ProjectName string                      `json:"projectName" yaml:"projectName"`
	SubType     metav1alpha1.ProjectSubType `json:"subType" yaml:"subType"`
}

// SubtypeProjectGetter get project with subtype
// alias of ProjectGetter for better experience
type SubtypeProjectGetter interface {
	Interface
	GetSubTypeProject(ctx context.Context, getOption GetProjectOption) (*metav1alpha1.Project, error)
}

// ProjectCreator create project api
type ProjectCreator interface {
	Interface
	CreateProject(ctx context.Context, project *metav1alpha1.Project) (*metav1alpha1.Project, error)
}

// ProjectDeleter create project api
type ProjectDeleter interface {
	Interface
	DeleteProject(ctx context.Context, project *metav1alpha1.Project) error
}
