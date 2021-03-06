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
	v1 "k8s.io/api/core/v1"
)

func (r *Repository) AddNamespaceRef(refs ...v1.ObjectReference) {
	for _, item := range refs {
		found := false
		for _, old := range r.Spec.NamespaceRefs {
			if old.Name == item.Name {
				found = true
				break
			}
		}

		if !found {
			r.Spec.NamespaceRefs = append(r.Spec.NamespaceRefs, item)
		}
	}
}

// Paginate return a pagination subset of repository list with specific page and page size
func (r *RepositoryList) Paginate(page int, pageSize int) *RepositoryList {
	length := len(r.Items)
	skip, end := common.Paginate(length, pageSize, page)

	newList := &RepositoryList{}
	newList.Items = r.Items[skip:end]
	newList.ListMeta.TotalItems = length

	return newList
}

// Filter takes a closure that returns true or false, if true, the repository should be present
func (r *RepositoryList) Filter(filter func(repository Repository) bool) *RepositoryList {
	if filter == nil {
		return r
	}

	newList := &RepositoryList{}
	for _, repository := range r.Items {
		if filter(repository) {
			newList.Items = append(newList.Items, repository)
		}
	}

	return newList
}
