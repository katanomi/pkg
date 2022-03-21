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
	"github.com/katanomi/pkg/common"
)

// Paginate return a pagination subset of artifact list with specific page and page size
func (r *ArtifactList) Paginate(page int, pageSize int) *ArtifactList {
	length := len(r.Items)
	skip, end := common.Paginate(length, pageSize, page)

	newList := &ArtifactList{}
	newList.Items = r.Items[skip:end]
	newList.ListMeta.TotalItems = length

	return newList
}

// Filter takes a closure that returns true or false, if true, the artifact should be present
func (r *ArtifactList) Filter(filter func(artifact Artifact) bool) *ArtifactList {
	if filter == nil {
		return r
	}

	newList := &ArtifactList{}
	for _, artifact := range r.Items {
		if filter(artifact) {
			newList.Items = append(newList.Items, artifact)
		}
	}

	return newList
}
