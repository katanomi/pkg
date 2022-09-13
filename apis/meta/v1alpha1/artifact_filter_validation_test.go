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
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestArtifactFilterRegexList_Validate(t *testing.T) {
	var (
		list ArtifactFilterRegexList
		g    = NewGomegaWithT(t)
	)

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
		filter ArtifactEnvFilter
		g      = NewGomegaWithT(t)
	)

	g.Expect(filter.Validate(field.NewPath("test"))).NotTo(BeEmpty())

	filter.Name = "test"
	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())

	filter.Regex = ArtifactFilterRegexList{
		"abc",
		"abc def",
		"^abc$",
		"abc+?",
		"abc*",
		"(abc)",
		"[abc]",
		`\d\w`,
	}

	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())
}

func TestArtifactTagFilter_Validate(t *testing.T) {
	var (
		filter ArtifactTagFilter
		g      = NewGomegaWithT(t)
	)

	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())

	filter.Regex = ArtifactFilterRegexList{
		"abc",
		"abc def",
		"^abc$",
		"abc+?",
		"abc*",
		"(abc)",
		"[abc]",
		`\d\w`,
	}

	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())

	filter.Regex = ArtifactFilterRegexList{
		"abc",
		"[",
	}
	g.Expect(filter.Validate(field.NewPath("test"))).To(HaveLen(1))
}

func TestArtifactLabelFilter_Validate(t *testing.T) {
	var (
		filter ArtifactLabelFilter
		g      = NewGomegaWithT(t)
	)

	g.Expect(filter.Validate(field.NewPath("test"))).NotTo(BeEmpty())

	filter.Name = "test"
	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())

	filter.Regex = ArtifactFilterRegexList{
		"abc",
		"abc def",
		"^abc$",
		"abc+?",
		"abc*",
		"(abc)",
		"[abc]",
		`\d\w`,
	}

	g.Expect(filter.Validate(field.NewPath("test"))).To(BeEmpty())

	filter.Name = ""
	filter.Regex = ArtifactFilterRegexList{
		"abc",
		"[",
	}
	g.Expect(filter.Validate(field.NewPath("test"))).To(HaveLen(2))
}
