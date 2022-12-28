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
	"testing"

	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestDeployRepositoryOption_Validate(t *testing.T) {
	tests := []struct {
		name             string
		DeployRepository string
		Required         bool
		wantErr          bool
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
			DeployRepository: "test.com/\f&^test",
			wantErr:          true,
		},
		{
			name:     "validate required repstory",
			Required: true,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DeployRepositoryOption{DeployRepository: tt.DeployRepository, Required: tt.Required}
			if gotErrs := m.Validate(field.NewPath("unittest")); tt.wantErr != (len(gotErrs) != 0) {
				t.Errorf("DeployRepositoryOption.Validate() = %v, want %v", gotErrs, tt.wantErr)
			}
		})
	}
}

func TestDeployRepositoryOption_Setup(t *testing.T) {
	t.Run("parse flag", func(t *testing.T) {
		g := NewGomegaWithT(t)
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
		g.Expect(obj.DeployRepository).To(Equal("registry.com"))
	})

	t.Run("setup get args", func(t *testing.T) {
		g := NewGomegaWithT(t)
		obj := struct {
			DeployRepositoryOption
		}{}

		args := []string{
			"--deploy-args", "tag1=1", "tag2=2",
		}
		err := RegisterSetup(&obj, context.Background(), nil, args)
		g.Expect(err).Should(Succeed(), "step succeed.")
		g.Expect(obj.DeployArgs).To(HaveLen(2), "get args succeed")
	})

}
