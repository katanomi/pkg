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

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func TestDeployApplicationResults_IsEmpty(t *testing.T) {
	var deployApp *DeployApplicationResults
	t.Run("is nil, should return true", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		deployApp = nil
		g.Expect(deployApp.IsEmpty()).To(gomega.BeTrue())
	})
	t.Run("only has application reference, should return true", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		deployApp = &DeployApplicationResults{ApplicationRef: &corev1.ObjectReference{}}
		g.Expect(deployApp.IsEmpty()).To(gomega.BeTrue())
	})
	t.Run("has empty data should return true", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		// only checks if there is an item, the contents are not checked
		deployApp = &DeployApplicationResults{
			ApplicationRef: &corev1.ObjectReference{},
			After: []DeployApplicationStatus{
				{},
			},
			Before: []DeployApplicationStatus{
				{},
			},
		}
		g.Expect(deployApp.IsEmpty()).To(gomega.BeTrue())
	})
	t.Run("has application reference before with data, should return false", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)

		deployApp = &DeployApplicationResults{
			ApplicationRef: &corev1.ObjectReference{
				Name: "abc",
			},
			Before: []DeployApplicationStatus{
				{
					Name: "abc",
				},
			},
		}
		g.Expect(deployApp.IsEmpty()).To(gomega.BeFalse())
	})
}
