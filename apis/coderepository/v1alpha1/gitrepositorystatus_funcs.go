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

package v1alpha1

import (
	"context"

	"github.com/katanomi/pkg/maps"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// GetValWithKey get variables with values for variable substitution using
// a cascading method from each different object using path as a path builder
// from top to bottom. It returns a set of string variables, array variables
// and object variables like the following:
//
// | Variable               |  Example value                              |
// |------------------------|---------------------------------------------|
// | <path>.<name>          | {"url":"","revision:"","commit":""}         |
// | <path>.<name>.revision | "refs/heads/main"                           |
// | <path>.<name>.commit   | "e683afce4427afc58ee8af4c53c030854616aff2"  |
// | <path>.<name>.url      | "e683afce4427afc58ee8af4c53c030854616aff2"  |
// etc..
func (git *GitRepositoryStatus) GetValWithKey(ctx context.Context, path *field.Path) (stringVals map[string]string, arrayVals map[string][]string, objectVals map[string]map[string]string) {
	stringVals = map[string]string{}
	arrayVals = map[string][]string{}
	objectVals = map[string]map[string]string{}
	if git == nil {
		// empty values returned here
		return
	}
	path = path.Child(git.Name)
	stringVals[path.Child("name").String()] = git.Name
	if git.BaseGitStatus != nil {
		stringVals = maps.MergeMap(stringVals, git.BaseGitStatus.GetValWithKey(ctx, path))
		// The following logic is to fulfil git repository contract
		// and it should be moved to the contract when a new way of doing this is
		// designed. For now it is quite manual
		// TODO: move to be implemented inside contract

		// Current implementation will add an empty root key inside
		// git base status  <path>.<name>
		// but we want to leave it to equate the contract properties
		// url, revision, commit
		// delete(stringVals, path.String())
		if git.BaseGitStatus.LastCommit != nil && git.BaseGitStatus.Revision != nil {
			delete(stringVals, path.String())
			stringVals = maps.MergeMap(stringVals, map[string]string{
				path.Child("url").String():      git.BaseGitStatus.URL,
				path.Child("revision").String(): git.BaseGitStatus.Revision.Raw,
				path.Child("commit").String():   git.BaseGitStatus.LastCommit.ID,
			})
			objectVals = maps.MergeMapMap(objectVals, map[string]map[string]string{
				path.String(): {
					"url":      git.BaseGitStatus.URL,
					"revision": git.BaseGitStatus.Revision.Raw,
					"commit":   git.BaseGitStatus.LastCommit.ID,
				},
			})
		}
	}
	return
}
