/*
Copyright 2021 The AlaudaDevops Authors.

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

package common

const (
	// Default number of pages
	DefaultPerPage = 20

	// Default page number
	DefaultPage = 1
)

// Paginate this function implements the paging of List, inputting the total number of data (usually
// the length of Slice), the number of each page, and the page number. The function guarantees
// that Begin and End do not exceed the total number of data.
//
// If input perPage is less than 1, it will be set as DefaultPerPage.
// If input page is less than 1, it will be set as DefaultPage.
//
// Returns the start and end position of the pagination.
func Paginate(total, perPage, page int) (begin, end int) {
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if page < 1 {
		page = DefaultPage
	}

	begin = (page - 1) * perPage
	end = begin + perPage
	if begin > total {
		begin = total
	}
	if end > total {
		end = total
	}

	return
}
