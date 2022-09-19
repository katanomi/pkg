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

import "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

func FilterCompletedTaskRun(list *v1beta1.TaskRunList) {
	items := list.Items
	for i := 0; i < len(items); i++ {
		if !IsCompletedTaskRun(items[i]) {
			items = append(items[:i], items[i+1:]...)
			i--
		}
	}
	list.Items = items
}

func IsCompletedTaskRun(taskRun v1beta1.TaskRun) bool {
	// When taskRun is Done the completionTime should be set
	// In case completionTime is nil cause crash we check if completionTime is nil here.
	return taskRun.IsDone() && taskRun.Status.CompletionTime != nil
}
