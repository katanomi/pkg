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

// Equal compares two branchSpec for equality.
func (b *BranchSpec) Equal(item BranchSpec) bool {
	return b.CodeInfo.Equal(item.CodeInfo) &&
		b.Author.Equal(item.Author) &&
		b.Issue.Equal(item.Issue)
}

// Equal compares two userSpec for equality.
func (a *UserSpec) Equal(item UserSpec) bool {
	return a.Id == item.Id &&
		a.Name == item.Name &&
		a.Email == item.Email
}

// Equal compares two codeinfos for equality, the IntegrationName is not used as the basis for judgment.
func (c *CodeInfo) Equal(item CodeInfo) bool {
	var host1, host2 string
	if c.Address != nil && c.Address.URL != nil {
		host1 = c.Address.URL.String()
	}

	if item.Address != nil && item.Address.URL != nil {
		host2 = item.Address.URL.String()
	}

	return c.Project == item.Project &&
		c.Repository == item.Repository &&
		c.Branch == item.Branch &&
		c.BaseBranch == item.BaseBranch &&
		host1 == host2
}

// Equal compares two issueinfo for equality.
func (c *IssueInfo) Equal(item IssueInfo) bool {
	return c.Type == item.Type && c.Id == item.Id
}
