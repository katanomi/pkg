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
	"fmt"
	"strconv"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func TestProjectAddNamespaceRef(t *testing.T) {
	cases := []struct {
		Refs     []v1.ObjectReference
		Expected int
	}{
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-1",
				},
			},
			Expected: 1,
		},
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-2",
				},
			},
			Expected: 2,
		},
		{
			Refs: []v1.ObjectReference{
				{
					Name: "test-1",
				},
				{
					Name: "test-2",
				},
				{
					Name: "test-1",
				},
			},
			Expected: 2,
		},
	}

	for i := range cases {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			g := NewGomegaWithT(t)
			p := &Project{}
			p.AddNamespaceRef(cases[i].Refs...)
			g.Expect(p.Spec.NamespaceRefs).To(HaveLen(cases[i].Expected))
		})
	}
}

func TestProjectList_Paginate(t *testing.T) {
	list := ProjectList{}
	for i := 0; i < 1000; i++ {
		project := Project{
			ObjectMeta: metav1.ObjectMeta{
				Name: strconv.Itoa(i),
			},
		}
		list.Items = append(list.Items, project)
	}

	g := NewGomegaWithT(t)

	for j := 1; j < 100; j++ {
		p := list.Paginate(j, 10)
		g.Expect(p.TotalItems).To(Equal(1000))
		g.Expect(len(p.Items)).To(Equal(10))
		g.Expect(p.Items[0]).To(Equal(list.Items[(j-1)*10]))
		g.Expect(p.Items[9]).To(Equal(list.Items[j*10-1]))
	}
}

func TestProjectList_Sort(t *testing.T) {
	cases := []struct {
		Names    []string
		Expected []string
	}{
		{
			Names:    []string{"f", "c", "b", "a", "e", "d"},
			Expected: []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			Names:    []string{"a", "1", "A"},
			Expected: []string{"1", "A", "a"},
		},
		{
			Names:    []string{"a", "aa", "1", "11"},
			Expected: []string{"1", "11", "a", "aa"},
		},
	}

	g := NewGomegaWithT(t)

	for i, item := range cases {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			list1 := &ProjectList{}
			for _, name := range item.Names {
				project := Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
					},
				}

				list1.Items = append(list1.Items, project)
			}

			list2 := &ProjectList{}
			for _, name := range item.Expected {
				project := Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
					},
				}

				list2.Items = append(list2.Items, project)
			}
			g.Expect(list1.Sort().Items).To(Equal(list2.Items))
		})
	}
}

func TestProjectList_Filter(t *testing.T) {
	cases := []struct {
		Names    []string
		Filter   string
		Expected []string
	}{
		{
			Names:    []string{"a", "ab", "abc"},
			Filter:   "a",
			Expected: []string{"a", "ab", "abc"},
		},
		{
			Names:  []string{"a", "ab", "abc"},
			Filter: "b",

			Expected: []string{"ab", "abc"},
		},
		{
			Names:    []string{"a", "ab", "abc"},
			Filter:   "c",
			Expected: []string{"abc"},
		},
		{
			Names:    []string{"a", "ab", "abc"},
			Filter:   "d",
			Expected: []string{},
		},
	}

	g := NewGomegaWithT(t)

	for i, item := range cases {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			list1 := &ProjectList{}
			for _, name := range item.Names {
				project := Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
					},
				}

				list1.Items = append(list1.Items, project)
			}
			list1 = list1.Filter(func(project Project) bool {
				return strings.Contains(project.Name, item.Filter)
			})

			list2 := &ProjectList{}
			for _, name := range item.Expected {
				project := Project{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
					},
				}

				list2.Items = append(list2.Items, project)
			}
			g.Expect(list1.Items).To(Equal(list2.Items))
		})
	}
}

func TestFilterProject(t *testing.T) {
	tests := []struct {
		project string
		include string
		exclude string
		result  bool
	}{
		{
			project: "abc",
			include: "def",
			exclude: "",
			result:  false,
		},
		{
			project: "abc",
			include: "abc",
			exclude: "",
			result:  true,
		},
		{
			project: "abc",
			include: "abc",
			exclude: "def",
			result:  true,
		},
		{
			project: "abc",
			include: "abc",
			exclude: "abc",
			result:  false,
		},
		{
			project: "abc",
			include: "def,123",
			exclude: "",
			result:  false,
		},
		{
			project: "abc",
			include: "abc,123",
			exclude: "",
			result:  true,
		},
		{
			project: "abc",
			include: "abc,123",
			exclude: "def,456",
			result:  true,
		},
		{
			project: "abc",
			include: "abc,123",
			exclude: "abc,456",
			result:  false,
		},
		{
			project: "abc,bcd",
			include: "abc",
			exclude: "abc,456",
			result:  false,
		},
		{
			project: "abc,bcd",
			include: "abc,123",
			exclude: "abc,456",
			result:  false,
		},
		{
			project: "abc,bcd",
			include: "abc,123",
			exclude: "456",
			result:  true,
		},
		{
			project: "abc,bcd",
			include: "abc",
			exclude: "",
			result:  true,
		},
		{
			project: "abc,bcd",
			include: "abc",
			exclude: "bcd",
			result:  true,
		},
	}

	g := NewGomegaWithT(t)

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			g.Expect(FilterProject(item.project, item.include, item.exclude)).To(Equal(item.result))
		})
	}
}
