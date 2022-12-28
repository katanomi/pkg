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

package options

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestDeployRepositoryOption_Validate(t *testing.T) {
	tests := []struct {
		name             string
		DeployRepository string
		wantErrs         field.ErrorList
	}{
		{
			name:             "validate success",
			DeployRepository: "http://test.com/test",
		},
		{
			name: "validate repository is empty",
		},
		{
			name:             "validate repository failed",
			DeployRepository: "test.com/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DeployRepositoryOption{DeployRepository: tt.DeployRepository}
			if gotErrs := m.Validate(field.NewPath("unittest")); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("DeployRepositoryOption.Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestDeployRepositoryOption_Setup(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.Background()
	base := field.NewPath("base")

	obj := struct {
		DeployRepositoryOption
	}{}
	args := []string{
		"--deploy-repository", "registry.com",
	}

	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)
	err := flagSet.Parse(args)
	g.Expect(err).Should(Succeed(), "parse flag succeed.")

	args = []string{
		"--deploy-args", "tag=1", "tag=2",
	}
	err = RegisterSetup(&obj, ctx, nil, args)
	g.Expect(err).Should(Succeed(), "step flag succeed.")
	g.Expect(obj.DeployRepository).To(Equal("registry.com"))
	g.Expect(obj.Validate(base)).To(HaveLen(0), "validate succeed")
}
