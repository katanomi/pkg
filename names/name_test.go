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

package names

import (
	// "context"
	"testing"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetGenerateName(t *testing.T) {

	table := map[string]struct {
		Object metav1.Object
		Result string
	}{
		"secret with simple name": {
			Object: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "abc",
				},
			},
			Result: "abc-",
		},
		"secret with generate name": {
			Object: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "xyz",
				},
			},
			Result: "xyz-",
		},
	}

	for i, tst := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := GetGenerateName(tst.Object)

			g.Expect(result).To(Equal(tst.Result))

		})
	}
}

func TestGenerateName(t *testing.T) {

	table := map[string]struct {
		Base        string
		Result      string
		ExpectedLen int
	}{
		"basic case": {
			Base:        "abc",
			Result:      "abc",
			ExpectedLen: 8,
		},
		"empty string case": {
			Base:        "",
			Result:      "",
			ExpectedLen: 5,
		},
		"full 63 characters case": {
			Base:        "abc4567890abc4567890abc4567890abc4567890abc4567890abc4567890abc",
			Result:      "abc4567890abc4567890abc4567890abc4567890abc4567890abc45678",
			ExpectedLen: 63,
		},
	}

	for i, tst := range table {
		t.Run(i, func(t *testing.T) {
			g := NewGomegaWithT(t)

			result := GenerateName(tst.Base)

			g.Expect(result).To(ContainSubstring(tst.Result))
			g.Expect(result).To(HaveLen(tst.ExpectedLen))
		})
	}
}
