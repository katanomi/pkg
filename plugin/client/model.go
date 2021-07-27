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

package client

type ListMeta struct {
	// Total number of items on the list. Used for pagination.
	TotalItems int `json:"totalItems"`
}

type ProjectList struct {
	ListMeta ListMeta   `json:"listMeta"`
	Items    []*Project `json:"items"`
	Errors   []error    `json:"errors"`
}

type Project struct {
	ApiVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ProjectMetadata `json:"metadata"`
	Spec       ProjectSpec     `json:"spec"`
}

type ProjectMetadata struct {
	Name string `json:"name"`
}

type ProjectSpec struct {
	Name         string   `json:"name"`
	BindProjects []string `json:"bindProjects"`
	Public       bool     `json:"public"`
	URL          string   `json:"url"`
	Attributes   map[string]interface{}
}

type Resource struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

type ResourceList []*Resource
