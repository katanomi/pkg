/*
Copyright 2022 The Katanomi Authors.

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

package tekton

import (
	"sort"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var _ sort.Interface = &TaskRunListSortByCompletion{}

// TaskRunListSortByCompletion represent tekton.TaskRunList but will
// sort by completion time
type TaskRunListSortByCompletion v1beta1.TaskRunList

// Len implement sort interface
func (s *TaskRunListSortByCompletion) Len() int {
	return len(s.Items)
}

// Less implement sort interface
func (s *TaskRunListSortByCompletion) Less(i, j int) bool {
	iTime := s.Items[i].Status.CompletionTime
	jTime := s.Items[j].Status.CompletionTime
	return jTime.Before(iTime)
}

// Swap implement sort interface
func (s *TaskRunListSortByCompletion) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
	return
}

// SortTaskRunByCompletion will sort TaskRunList by completion time
func SortTaskRunByCompletion(list *v1beta1.TaskRunList) {
	result := TaskRunListSortByCompletion(*list)
	sort.Sort(&result)
	list.Items = result.Items
	return
}
