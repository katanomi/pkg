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

package scheme

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestSchemeContext(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := context.TODO()

	clt := Scheme(ctx)
	g.Expect(clt).To(BeNil())

	ctx = WithScheme(ctx, scheme.Scheme)
	g.Expect(Scheme(ctx)).To(Equal(scheme.Scheme))
}
