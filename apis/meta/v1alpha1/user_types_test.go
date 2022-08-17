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
	"context"

	. "github.com/katanomi/pkg/testing"
	rbacv1 "k8s.io/api/rbac/v1"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/onsi/gomega"
)

var _ = Describe("RBACSubjectValGetter.GetValWithKey", func() {
	var (
		ctx      context.Context
		path     *field.Path
		subject  *rbacv1.Subject
		values   map[string]string
		expected map[string]string
	)
	BeforeEach(func() {
		ctx = context.TODO()
		path = field.NewPath("user")
		subject = &rbacv1.Subject{}
		expected = map[string]string{}
	})
	JustBeforeEach(func() {
		values = RBACSubjectValGetter(subject)(ctx, path)
	})
	Context("rbac.Subject with all variables", func() {
		BeforeEach(func() {
			Expect(LoadYAML("testdata/user_types_vars.all.yaml", subject)).To(Succeed())
			Expect(LoadYAML("testdata/user_types_vars.all.golden.yaml", &expected)).To(Succeed())
			Expect(expected).ToNot(BeEmpty())
		})
		It("should return the same amount of data", func() {
			Expect(values).To(Equal(expected))
		})
	})
})
