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

package v1alpha1

// Pager describe the paging params
type Pager struct {
	// ItemsPerPage desired number of items to be returned in each page
	ItemsPerPage int `json:"itemsPerPage"`

	// Page desired to be returned
	Page int `json:"page"`
}

// GetPageLimit get the limit returned by a single page
func (p *Pager) GetPageLimit() int {
	if p.ItemsPerPage == 0 {
		p.ItemsPerPage = 20
	}
	return p.ItemsPerPage
}

// GetOffset get the offset for next query
func (p *Pager) GetOffset() int {
	return p.GetPageLimit() * (p.GetPage() - 1)
}

// GetPage get the current page
func (p *Pager) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
