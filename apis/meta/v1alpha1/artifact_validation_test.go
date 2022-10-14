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
	"context"
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/gomega"
)

func TestArtifactParameterSpec_Validate(t *testing.T) {

	tests := []struct {
		init     func(a *ArtifactParameterSpec)
		evaluate func(g *GomegaWithT, errs field.ErrorList)
	}{
		{
			init: func(a *ArtifactParameterSpec) {
			},
			evaluate: func(g *GomegaWithT, errs field.ErrorList) {
				g.Expect(errs[0].Error()).To(ContainSubstring("Required"))
			},
		},
		{
			init: func(a *ArtifactParameterSpec) {
				a.URI = "127.0.0.1"
			},
			evaluate: func(g *GomegaWithT, errs field.ErrorList) {
				g.Expect(errs).To(BeEmpty())
			},
		},
		{
			init: func(a *ArtifactParameterSpec) {
				a.URI = "127.0.0.1"
				a.IntegrationClassName = "harbor"
			},
			evaluate: func(g *GomegaWithT, errs field.ErrorList) {
				g.Expect(errs).To(BeEmpty())
			},
		},
		{
			init: func(a *ArtifactParameterSpec) {
				a.URI = "127.0.0.1/repo/dest"
				a.IntegrationClassName = ""
			},
			evaluate: func(g *GomegaWithT, errs field.ErrorList) {
				g.Expect(errs).To(BeEmpty())
			},
		},
		{
			init: func(a *ArtifactParameterSpec) {
				a.URI = " 127.0.0.1"
				a.IntegrationClassName = "docker-registry"
			},
			evaluate: func(g *GomegaWithT, errs field.ErrorList) {
				g.Expect(errs[0].Error()).To(ContainSubstring("invalid"))
			},
		},
	}

	g := NewGomegaWithT(t)
	for i, item := range tests {
		a := &ArtifactParameterSpec{}
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			item.init(a)
			err := a.Validate(context.TODO(), field.NewPath("test"))
			item.evaluate(g, err)
		})
	}
}
