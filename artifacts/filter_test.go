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

package artifacts

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

func TestLabelFilter(t *testing.T) {

	var data = []struct {
		desc   string
		labels map[string]string
		filter v1alpha1.ArtifactLabelFilter

		res FilterResult
		err bool
	}{
		{
			desc: "has filter, but labels is empty",

			filter: v1alpha1.ArtifactLabelFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			labels: map[string]string{},

			res: FailFilter,
			err: false,
		},
		{
			desc: "has filter, but not contains target label",

			filter: v1alpha1.ArtifactLabelFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			labels: map[string]string{
				"qa": "sign-off",
			},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, contains target label, but not matched values",
			filter: v1alpha1.ArtifactLabelFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			labels: map[string]string{
				"qa":     "sign-off",
				"branch": "develop",
			},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, contains target label, matched values in one regex",
			filter: v1alpha1.ArtifactLabelFilter{Name: "branch", Regex: []string{"release-.*"}},
			labels: map[string]string{
				"qa":     "sign-off",
				"branch": "release-1.0",
			},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has filter, contains target label, matched values in more than one regex",
			filter: v1alpha1.ArtifactLabelFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			labels: map[string]string{
				"qa":     "sign-off",
				"branch": "release-1.0",
			},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has no filter, labels is empty",
			filter: v1alpha1.ArtifactLabelFilter{},
			labels: map[string]string{},

			res: NoFilter,
			err: false,
		},
		{
			desc: "has no filter, labels is not empty",
			labels: map[string]string{
				"qa":     "sign-off",
				"branch": "release-1.0",
			},

			res: NoFilter,
			err: false,
		},
	}

	a := v1alpha1.Artifact{Spec: v1alpha1.ArtifactSpec{Properties: &runtime.RawExtension{}}}
	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {

			a.Spec.Properties.Raw = makeProperties(item.labels, []string{}, []string{})
			f := NewLabelFilter(item.filter)
			actualRes, message, err := f.Filter(context.TODO(), a)
			if actualRes == FailFilter {
				fmt.Println(message)
			}
			if item.err && err == nil {
				t.Errorf("should filter with error, but actual err is nil")
			}

			if !item.err && err != nil {
				t.Errorf("should filter with not error, but acutal err is not nil: %s", err.Error())
			}

			if actualRes != item.res {
				t.Errorf("expect filter result is %s, acutal is %s", item.res, actualRes)
			}
		})

	}
}

func makeProperties(labels map[string]string, envs []string, tags []string) []byte {

	obj := map[string]interface{}{}

	objTags := []map[string]string{}
	if len(tags) > 0 {
		for _, tag := range tags {
			objTags = append(objTags, map[string]string{"name": tag})
		}
	}
	obj["tags"] = objTags

	obj["extra_attrs"] = map[string]interface{}{
		"config": map[string]interface{}{
			"Labels": labels,
			"Env":    envs,
		},
	}

	bts, _ := json.Marshal(obj)

	return bts
}

func TestTagFilter(t *testing.T) {
	var data = []struct {
		desc   string
		tags   []string
		filter v1alpha1.ArtifactTagFilter

		res FilterResult
		err bool
	}{
		{
			desc: "has filter, but tags is empty",

			filter: v1alpha1.ArtifactTagFilter{Regex: []string{"^master$", "release-.*"}},
			tags:   []string{},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, but not matched values",
			filter: v1alpha1.ArtifactTagFilter{Regex: []string{"^master$", "release-.*"}},
			tags:   []string{"develop", "feature"},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, matched values in one regex",
			filter: v1alpha1.ArtifactTagFilter{Regex: []string{"release-.*"}},
			tags:   []string{"release-3.5"},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has filter, matched values in more than one regex",
			filter: v1alpha1.ArtifactTagFilter{Regex: []string{"^master$", "release-.*"}},
			tags:   []string{"release-3.5", "develop"},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has no filter, tags is empty",
			filter: v1alpha1.ArtifactTagFilter{},
			tags:   []string{},

			res: NoFilter,
			err: false,
		},
		{
			desc: "has no filter, labels is not empty",
			tags: []string{"release-3.5", "develop"},

			res: NoFilter,
			err: false,
		},
	}

	a := v1alpha1.Artifact{Spec: v1alpha1.ArtifactSpec{Properties: &runtime.RawExtension{}}}
	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {

			a.Spec.Properties.Raw = makeProperties(map[string]string{}, []string{}, item.tags)
			f := NewTagFilter(item.filter)
			actualRes, message, err := f.Filter(context.TODO(), a)
			if actualRes == FailFilter {
				fmt.Println(message)
			}

			if item.err && err == nil {
				t.Errorf("should filter with error, but actual err is nil")
			}

			if !item.err && err != nil {
				t.Errorf("should filter with not error, but acutal err is not nil: %s", err.Error())
			}

			if actualRes != item.res {
				t.Errorf("expect filter result is %s, acutal is %s", item.res, actualRes)
			}
		})
	}
}

func TestEnvFilter(t *testing.T) {
	var data = []struct {
		desc   string
		envs   []string
		filter v1alpha1.ArtifactEnvFilter

		res FilterResult
		err bool
	}{
		{
			desc: "has filter, but envs is empty",

			filter: v1alpha1.ArtifactEnvFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			envs:   []string{},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, but not contains name",
			filter: v1alpha1.ArtifactEnvFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			envs:   []string{"GOBIN=/go/bin", "PATH=/usr/bin"},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, but not matched values",
			filter: v1alpha1.ArtifactEnvFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			envs:   []string{"GOBIN=/go/bin", "PATH=/usr/bin", "branch=develop"},

			res: FailFilter,
			err: false,
		},
		{
			desc:   "has filter, matched values in one regex",
			filter: v1alpha1.ArtifactEnvFilter{Name: "branch", Regex: []string{"release-.*"}},
			envs:   []string{"GOBIN=/go/bin", "PATH=/usr/bin", "branch=release-1.0"},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has filter, matched values in more than one regex",
			filter: v1alpha1.ArtifactEnvFilter{Name: "branch", Regex: []string{"^master$", "release-.*"}},
			envs:   []string{"GOBIN=/go/bin", "PATH=/usr/bin", "branch=release-1.0"},

			res: PassFilter,
			err: false,
		},
		{
			desc:   "has no filter, envs is empty",
			filter: v1alpha1.ArtifactEnvFilter{},
			envs:   []string{},

			res: NoFilter,
			err: false,
		},
		{
			desc: "has no filter, envs is not empty",
			envs: []string{"GOBIN=/go/bin", "PATH=/usr/bin", "branch=release-1.0"},

			res: NoFilter,
			err: false,
		},
	}

	a := v1alpha1.Artifact{Spec: v1alpha1.ArtifactSpec{Properties: &runtime.RawExtension{}}}
	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {

			a.Spec.Properties.Raw = makeProperties(map[string]string{}, item.envs, []string{})
			f := NewEnvFilter(item.filter)
			actualRes, message, err := f.Filter(context.TODO(), a)
			if actualRes == FailFilter {
				fmt.Println(message)
			}

			if item.err && err == nil {
				t.Errorf("should filter with error, but actual err is nil")
			}

			if !item.err && err != nil {
				t.Errorf("should filter with not error, but acutal err is not nil: %s", err.Error())
			}

			if actualRes != item.res {
				t.Errorf("expect filter result is %s, acutal is %s", item.res, actualRes)
			}
		})
	}
}

func TestAllFilter(t *testing.T) {
	var data = []struct {
		desc    string
		filters []v1alpha1.ArtifactFilter

		envs   []string
		labels map[string]string
		tags   []string

		res FilterResult
		err bool
	}{
		{
			desc:    "has no filter",
			filters: []v1alpha1.ArtifactFilter{},

			res: NoFilter,
			err: false,
		},
		{
			desc: "has one filter and matched",
			filters: []v1alpha1.ArtifactFilter{
				{
					Tags: []v1alpha1.ArtifactTagFilter{
						{
							Regex: []string{"^master$", "release-.*"},
						},
					},
				},
			},
			tags: []string{"release-1.0"},

			res: PassFilter,
			err: false,
		},
		{
			desc: "has more than one filter and matched",
			filters: []v1alpha1.ArtifactFilter{
				{
					Tags: []v1alpha1.ArtifactTagFilter{
						{
							Regex: []string{"^master$", "release-.*"},
						},
					},
					Labels: []v1alpha1.ArtifactLabelFilter{
						{
							Name:  "branch",
							Regex: []string{"^master$", "release-.*"},
						},
						{
							Name:  "qa",
							Regex: []string{"^sign-off$"},
						},
					},
				},
			},
			tags: []string{"release-1.0"},
			labels: map[string]string{
				"branch": "release-1.0",
				"qa":     "sign-off",
			},

			res: PassFilter,
			err: false,
		},
		{
			desc: "has more than one filter and one matched",
			filters: []v1alpha1.ArtifactFilter{
				{
					Tags: []v1alpha1.ArtifactTagFilter{
						{
							Regex: []string{"^master$", "release-.*"},
						},
					},
					Labels: []v1alpha1.ArtifactLabelFilter{
						{
							Name:  "branch",
							Regex: []string{"^master$", "release-.*"},
						},
						{
							Name:  "qa",
							Regex: []string{"^sign-off$"},
						},
					},
				},
				{
					Tags: []v1alpha1.ArtifactTagFilter{
						{
							Regex: []string{"^v1.0$"},
						},
					},
				},
			},
			tags: []string{"v1.0"},
			labels: map[string]string{
				"branch": "develop",
				"qa":     "testing",
			},

			res: FailFilter,
			err: false,
		},
	}

	a := v1alpha1.Artifact{Spec: v1alpha1.ArtifactSpec{Properties: &runtime.RawExtension{}}}
	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {

			a.Spec.Properties.Raw = makeProperties(item.labels, item.envs, item.tags)
			f := NewAllFilter(item.filters)
			actualRes, message, err := f.Filter(context.TODO(), a)
			if actualRes == FailFilter {
				fmt.Println(message)
			}

			if item.err && err == nil {
				t.Errorf("should filter with error, but actual err is nil")
			}

			if !item.err && err != nil {
				t.Errorf("should filter with not error, but acutal err is not nil: %s", err.Error())
			}

			if actualRes != item.res {
				t.Errorf("expect filter result is %s, acutal is %s", item.res, actualRes)
			}
		})
	}
}

func TestArtifactFilterSet(t *testing.T) {

	var artifact = v1alpha1.Artifact{Spec: v1alpha1.ArtifactSpec{Properties: &runtime.RawExtension{}}}
	artifact.Spec.Properties.Raw = makeProperties(
		map[string]string{
			"branch": "release-1.0",
			"qa":     "sign-off",
		},
		[]string{"branch=release-1.0", "qa=sign-off"},
		[]string{"v1.0"},
	)

	var data = []struct {
		desc string
		yaml string

		res FilterResult
		err bool
	}{
		{
			desc: "has no filter, it should pass",
			yaml: "",

			res: NoFilter,
			err: false,
		},
		{
			desc: "has only one envs and matched",
			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^master$", "release-.*" ]
            `,
			res: PassFilter,
			err: false,
		},
		{
			desc: "has only one envs but not match",
			yaml: `
          any:
          - envs:
            - name: branch
              regex: [ "^develop$"]
            `,
			res: FailFilter,
			err: false,
		},
		{
			desc: "has only envs and more than one in array and matched",

			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^master$", "release-.*" ]
              - name: qa
                regex: [ "^sign-off$"]
            `,
			res: PassFilter,
			err: false,
		},
		{
			desc: "has only envs and more than one in array but not matched",
			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^develop$"]
              - name: qa
                regex: [ "^sign-off$"]
            `,
			res: FailFilter,
			err: false,
		},
		{
			desc: "has envs and labels, more than one in array and not matched",
			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^develop$"]
              - name: qa
                regex: [ "^sign-off$"]
              labels:
              - name: branch
                regex: [ "^develop$"]
            `,
			res: FailFilter,
			err: false,
		},
		{
			desc: "has envs and labels, more than one in array and matched",
			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^master$", "release-.*"]
              - name: qa
                regex: [ "^sign-off$"]
              labels:
              - name: branch
                regex: [ "^release-.*$"]
            `,
			res: PassFilter,
			err: false,
		},
		{
			desc: "has envs and labels in any and matched",
			yaml: `
            any:
            - envs:
              - name: branch
                regex: [ "^master$"]
              - name: qa
                regex: [ "^sign-off$"]
            - labels:
              - name: branch
                regex: [ "^release-.*$"]
            `,
			res: PassFilter,
			err: false,
		},
		{
			desc: "has tags and envs and labels in any and matched",
			yaml: `
            any:
            - tags:
                - regex:
                    - "^v[1-9].*$"
              labels:
                - name: branch
                  regex:
                    - "^master$"
                    - "^release-.*$"
              envs:
                - name: branch
                  regex:
                    - "^master$"
                    - "^release-.*$"
            `,
			res: PassFilter,
			err: false,
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			set := v1alpha1.ArtifactFilterSet{}
			err := yaml.Unmarshal([]byte(item.yaml), &set)
			if err != nil {
				t.Errorf("error unmarshal filter set yaml: %s", err.Error())
			}
			actualRes, message, err := NewFilter(set).Filter(context.TODO(), artifact)
			if actualRes == FailFilter {
				fmt.Println(message)
			}

			if item.err && err == nil {
				t.Errorf("should filter with error, but actual err is nil")
			}

			if !item.err && err != nil {
				t.Errorf("should filter with not error, but acutal err is not nil: %s", err.Error())
			}

			if actualRes != item.res {
				t.Errorf("expect filter result is %s, acutal is %s", item.res, actualRes)
			}
		})
	}
}
