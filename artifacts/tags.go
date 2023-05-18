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

package artifacts

import (
	"context"
)

// FindMatchedTagByPrefixOptions options for FindMatchedTagByPrefix function
type FindMatchedTagByPrefixOptions struct {
	// MinNumberOfCharacterMatches stores the minimum allowed number of characters to be matched
	// in order to consider a valid candidate
	// defaults to 0 (no minimal limit)
	MinNumberOfCharacterMatches int
}

// FindMatchedTagByPrefix selects the best matched tag using a prefix, for example
// current tag = v1.8.0
// candidates are v1.8.1, v1.9.0, v2.0.0
// will return v1.8.1 because the most number for characters were matched
func FindMatchedTagByPrefix(ctx context.Context, currentTag string, opts FindMatchedTagByPrefixOptions, candidates ...string) (tag string, ok bool) {
	strongestIndex := -1
	topMatch := 0
	for idx, candidate := range candidates {
		matched := 0
		for ; matched < len(currentTag) && matched < len(candidate) && currentTag[matched] == candidate[matched]; matched++ {
			// no-op
		}
		if matched > topMatch && matched >= opts.MinNumberOfCharacterMatches {
			topMatch = matched
			strongestIndex = idx
		}
	}
	if strongestIndex >= 0 {
		tag = candidates[strongestIndex]
		ok = true
	}
	return
}
