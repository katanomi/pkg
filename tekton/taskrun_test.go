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
	"time"

	ktesting "github.com/katanomi/pkg/testing"

	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("taskRunSortDESCByCompleted", func() {
	var (
		taskRunEarly, taskRun, taskRunLate pipelinev1beta1.TaskRun

		taskRunList pipelinev1beta1.TaskRunList

		taskRunConstruct = func(name string, time metav1.Time) pipelinev1beta1.TaskRun {
			return pipelinev1beta1.TaskRun{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Status: pipelinev1beta1.TaskRunStatus{
					TaskRunStatusFields: pipelinev1beta1.TaskRunStatusFields{
						CompletionTime: &time,
					},
				},
			}
		}
	)

	BeforeEach(func() {
		completeTimeEarly := time.Now()
		completeTime := time.Now().Add(time.Minute)
		completeTimeLate := time.Now().Add(2 * time.Minute)

		taskRunEarly = taskRunConstruct("early", metav1.Time{Time: completeTimeEarly})
		taskRun = taskRunConstruct("now", metav1.Time{Time: completeTime})
		taskRunLate = taskRunConstruct("late", metav1.Time{Time: completeTimeLate})

		taskRunList.Items = append(taskRunList.Items, taskRunEarly)
		taskRunList.Items = append(taskRunList.Items, taskRunLate)
		taskRunList.Items = append(taskRunList.Items, taskRun)

	})

	Context("should return taskRun list", func() {
		It("by completedTime desc", func() {
			SortTaskRunByCompletion(&taskRunList)
			Expect(taskRunList.Items[0].Name).To(Equal("late"))
			Expect(taskRunList.Items[1].Name).To(Equal("now"))
			Expect(taskRunList.Items[2].Name).To(Equal("early"))
		})
	})

})

var _ = Describe("FilterCompletedTaskRun", func() {

	var (
		taskRunList *pipelinev1beta1.TaskRunList
		expectList  *pipelinev1beta1.TaskRunList
		g           *WithT
	)

	BeforeEach(func() {
		g = NewWithT(GinkgoT())
	})

	Context("no TaskRun provided", func() {
		BeforeEach(func() {
			taskRunList = &pipelinev1beta1.TaskRunList{}
			expectList = &pipelinev1beta1.TaskRunList{}
			FilterCompletedTaskRun(taskRunList)
		})
		It("should return no TaskRun", func() {
			Expect(taskRunList).To(Equal(expectList))
		})
	})

	Context("no TaskRun is completed", func() {
		BeforeEach(func() {
			taskRunList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("1")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("2")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("3")).(*pipelinev1beta1.TaskRun),
				},
			}
			expectList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{},
			}
			FilterCompletedTaskRun(taskRunList)
		})
		It("should return no TaskRun", func() {
			Expect(taskRunList).To(Equal(expectList))
		})
	})

	Context("first TaskRun is completed", func() {
		BeforeEach(func() {
			taskRunList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("1")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("2")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("3")).(*pipelinev1beta1.TaskRun),
				},
			}
			expectList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("1")).(*pipelinev1beta1.TaskRun),
				}}
			FilterCompletedTaskRun(taskRunList)
		})
		It("should return TaskRun 1", func() {
			Expect(taskRunList).To(Equal(expectList))
		})
	})

	Context("middle of TaskRun is completed", func() {
		BeforeEach(func() {
			taskRunList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("1")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("2")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("3")).(*pipelinev1beta1.TaskRun),
				},
			}
			expectList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("2")).(*pipelinev1beta1.TaskRun),
				}}
			FilterCompletedTaskRun(taskRunList)
		})
		It("should return TaskRun 2", func() {
			Expect(taskRunList).To(Equal(expectList))

		})
	})

	Context("latest TaskRun is completed", func() {
		BeforeEach(func() {
			taskRunList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("1")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Uncompleted.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("2")).(*pipelinev1beta1.TaskRun),
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("3")).(*pipelinev1beta1.TaskRun),
				},
			}
			expectList = &pipelinev1beta1.TaskRunList{
				Items: []pipelinev1beta1.TaskRun{
					*ktesting.LoadObjectOrDie(g, "testdata/TaskRun.Completed.yaml", &pipelinev1beta1.TaskRun{}, ktesting.SetName("3")).(*pipelinev1beta1.TaskRun),
				}}
			FilterCompletedTaskRun(taskRunList)
		})
		It("should return TaskRun 3", func() {
			Expect(taskRunList).To(Equal(expectList))
		})
	})
})

var _ = Describe("GetPipelineRunOwner", func() {

	type result struct {
		taskrun pipelinev1beta1.TaskRun
		exist   bool
		prName  string
	}

	var (
		withOwnerTR, noOwnerTR pipelinev1beta1.TaskRun
		expects                []result
	)

	BeforeEach(func() {
		Expect(ktesting.LoadYAML("testdata/TaskRun.WithOwner.yaml", &withOwnerTR)).To(Succeed())
		Expect(ktesting.LoadYAML("testdata/TaskRun.Completed.yaml", &noOwnerTR)).To(Succeed())
		expects = []result{
			{
				withOwnerTR,
				true,
				"complete",
			},
			{
				noOwnerTR,
				false,
				"",
			},
		}
	})

	Context("return owner if exist", func() {
		It("return owner as expect", func() {
			for _, r := range expects {
				exist, owner := GetPipelineRunOwner(r.taskrun)
				Expect(exist).To(Equal(r.exist))
				Expect(owner.Name).To(Equal(r.prName))
			}
		})
	})
})
