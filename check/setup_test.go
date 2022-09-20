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

package check

import (
	"context"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestInstalledTekton(t *testing.T) {
	testCases := []struct {
		clt  *fake.ClientBuilder
		want bool
	}{
		{
			clt:  fake.NewClientBuilder(),
			want: false,
		},
	}
	g := NewGomegaWithT(t)
	for _, tt := range testCases {
		client := tt.clt.Build()
		result := InstalledTekton(context.TODO(), client)
		g.Expect(result).To(Equal(tt.want))
	}
}
