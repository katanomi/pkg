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

import (
	"testing"
)

func TestRestoreEscapedCharacters(t *testing.T) {
	// key: string. value: except string
	caseMap := make(map[string]string)
	caseMap[`.`] = `.`
	caseMap[`\\`] = `\\`
	caseMap[`\\.`] = `\.`
	caseMap[`\.`] = `\.`
	caseMap[`cpaas\\.io`] = `cpaas\.io`
	caseMap[`cpaasio\\.`] = `cpaasio\.`
	caseMap[`/spec/replicationPolicies/0/namespaceFilter/filter/exact/$(metadata.labels.cpaas\\.io~1inner-namespace)`] = `/spec/replicationPolicies/0/namespaceFilter/filter/exact/$(metadata.labels.cpaas\.io~1inner-namespace)`
	caseMap[""] = ""
	caseMap[`..`] = `..`
	caseMap[`\\\.`] = `\\.`
	caseMap[`\\.abc`] = `\.abc`
	caseMap[`cpaas\\.iocpaas\\.io`] = `cpaas\.iocpaas\.io`

	for key, value := range caseMap {
		restore := restoreEscapedCharacters(key)

		if restore != value {
			t.Errorf("reduction should be %s, but got %s", value, restore)
		}
	}

}
