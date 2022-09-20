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

package v1alpha1

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestArtifactFilterRegexList_Validate(t *testing.T) {
	var (
		list ArtifactFilterRegexList
		g    = NewGomegaWithT(t)
	)

	g.Expect(list.Validate(field.NewPath("test"))).To(HaveLen(1))

	list = ArtifactFilterRegexList{
		"abc",
		"abc def",
		"^abc$",
		"abc+?",
		"abc*",
		"(abc)",
		"[abc]",
		`\d\w`,
	}

	g.Expect(list.Validate(field.NewPath("test"))).To(BeEmpty())

	list = ArtifactFilterRegexList{
		"abc",
		"[",
		"(",
	}
	g.Expect(list.Validate(field.NewPath("test"))).To(HaveLen(2))
}

func TestArtifactEnvFilter_Validate(t *testing.T) {
	var (
		g     = NewGomegaWithT(t)
		tests = []struct {
			name     string
			regex    ArtifactFilterRegexList
			evaluate func(g *GomegaWithT, errs field.ErrorList)
		}{
			{
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(2))
				},
			},
			{
				name: "test",
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				regex: ArtifactFilterRegexList{
					"abc",
					"abc def",
					"^abc$",
					"abc+?",
					"abc*",
					"(abc)",
					"[abc]",
					`\d\w`,
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				name: "test",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(BeEmpty())
				},
			},
			{
				name: "TEST",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(BeEmpty())
				},
			},
			{
				name: "?!#",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				name: "test",
				regex: ArtifactFilterRegexList{
					"abc",
					"[",
					"(",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(2))
				},
			},
			{
				name: "?!#",
				regex: ArtifactFilterRegexList{
					"abc",
					"[",
					"(",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(3))
				},
			},
		}
	)

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			filter := ArtifactEnvFilter{
				Name:  item.name,
				Regex: item.regex,
			}

			errs := filter.Validate(field.NewPath("test"))
			item.evaluate(g, errs)
		})
	}
}

func TestArtifactTagFilter_Validate(t *testing.T) {
	var (
		g     = NewGomegaWithT(t)
		tests = []struct {
			regex    ArtifactFilterRegexList
			evaluate func(g *GomegaWithT, errs field.ErrorList)
		}{
			{
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				regex: ArtifactFilterRegexList{
					"abc",
					"abc def",
					"^abc$",
					"abc+?",
					"abc*",
					"(abc)",
					"[abc]",
					`\d\w`,
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(BeEmpty())
				},
			},
			{
				regex: ArtifactFilterRegexList{
					"abc",
					"[",
					"(",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(2))
				},
			},
		}
	)

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			filter := ArtifactTagFilter{
				Regex: item.regex,
			}

			errs := filter.Validate(field.NewPath("test"))
			item.evaluate(g, errs)
		})
	}
}

func TestArtifactLabelFilter_Validate(t *testing.T) {
	var (
		g     = NewGomegaWithT(t)
		tests = []struct {
			name     string
			regex    ArtifactFilterRegexList
			evaluate func(g *GomegaWithT, errs field.ErrorList)
		}{
			{
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(2))
				},
			},
			{
				name: "test",
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				regex: ArtifactFilterRegexList{
					"abc",
					"abc def",
					"^abc$",
					"abc+?",
					"abc*",
					"(abc)",
					"[abc]",
					`\d\w`,
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				name: "test",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(BeEmpty())
				},
			},
			{
				name: "TEST",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(BeEmpty())
				},
			},
			{
				name: "?!#",
				regex: ArtifactFilterRegexList{
					"abc",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(1))
				},
			},
			{
				name: "test",
				regex: ArtifactFilterRegexList{
					"abc",
					"[",
					"(",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(2))
				},
			},
			{
				name: "?!#",
				regex: ArtifactFilterRegexList{
					"abc",
					"[",
					"(",
				},
				evaluate: func(g *GomegaWithT, errs field.ErrorList) {
					g.Expect(errs).To(HaveLen(3))
				},
			},
		}
	)

	for i, item := range tests {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			filter := ArtifactLabelFilter{
				Name:  item.name,
				Regex: item.regex,
			}

			errs := filter.Validate(field.NewPath("test"))
			item.evaluate(g, errs)
		})
	}
}

func TestArtifactFilterSet_Validate(t *testing.T) {
	var (
		g      = NewGomegaWithT(t)
		path   = field.NewPath("test")
		filter ArtifactFilterSet
	)

	filter.Any = nil
	filter.All = []ArtifactFilter{
		{},
	}
	g.Expect(filter.Validate(path)).To(BeEmpty())

	filter.All = nil
	filter.Any = []ArtifactFilter{
		{},
	}
	g.Expect(filter.Validate(path)).To(BeEmpty())

	filter.All = []ArtifactFilter{
		{},
	}
	filter.Any = []ArtifactFilter{
		{},
	}
	g.Expect(filter.Validate(path)).To(HaveLen(1))
}
