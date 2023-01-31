/*
Copyright 2023 The Katanomi Authors.

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

	. "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func TestNamedDeployApplicationResult_IsSameResult(t *testing.T) {
	t.Run("is the same name, content is different, should be true", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		original := NamedDeployApplicationResult{Name: "abc"}
		g.Expect(original.IsSameResult(NamedDeployApplicationResult{Name: "abc", DeployApplicationResults: DeployApplicationResults{ApplicationRef: &corev1.ObjectReference{}}})).To(gomega.BeTrue())
	})
	t.Run("different name, should be false", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		original := NamedDeployApplicationResult{Name: "abc"}
		g.Expect(original.IsSameResult(NamedDeployApplicationResult{Name: "def"})).To(gomega.BeFalse())
	})
}

func TestNamedDeployApplicationResults_DeepCopy(t *testing.T) {
	t.Run("NamedDeployApplicationResults.DeepCopy", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		obj := NamedDeployApplicationResults{}
		g.Expect(LoadYAML("testdata/deploy_application_status.DeepCopy.yaml", &obj)).To(gomega.Succeed())

		g.Expect(obj.DeepCopy()).To(gomega.Equal(obj))

		obj = nil
		g.Expect(obj.DeepCopy()).To(gomega.BeNil())
	})
	t.Run("*NamedDeployApplicationResult.DeepCopy", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		obj := &NamedDeployApplicationResult{}
		g.Expect(obj.DeepCopy()).To(gomega.Equal(obj))

		obj = nil
		g.Expect(obj.DeepCopy()).To(gomega.BeNil())
	})
	t.Run("*DeployApplicationResults.DeepCopy", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		obj := &DeployApplicationResults{}
		g.Expect(obj.DeepCopy()).To(gomega.Equal(obj))

		obj = nil
		g.Expect(obj.DeepCopy()).To(gomega.BeNil())
	})
	t.Run("*DeployApplicationStatus.DeepCopy", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		obj := &DeployApplicationStatus{}
		g.Expect(obj.DeepCopy()).To(gomega.Equal(obj))

		obj = nil
		g.Expect(obj.DeepCopy()).To(gomega.BeNil())
	})
}
