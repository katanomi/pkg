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
	"errors"
	"io"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/katanomi/pkg/apis/codequality/v1alpha1"
	"github.com/katanomi/pkg/command/qualitygate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

func TestCodeLinterOption_AddFlags(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := struct {
		CodeLinterOption
	}{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(&obj, flagSet)

	err := flagSet.Parse([]string{
		"--enable-quality-gate",
	})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.QualityGate).To(Equal(true))
}

func TestCodeLinterOption_Setup(t *testing.T) {
	g := NewGomegaWithT(t)
	obj := CodeLinterOption{}
	err := obj.Setup(context.Background(), nil, []string{"--quality-gate-rules", "a=b", "c=d"})
	g.Expect(err).Should(Succeed())
	g.Expect(obj.Result).To(Equal("Succeded"))
	g.Expect(obj.Issues).NotTo(BeNil())
	g.Expect(obj.QualityGateRules).To(Equal(map[string]string{
		"a": "b",
		"c": "d",
	}))
}

var _ = Describe("Test.CodeLinterOption.Validate", func() {
	var (
		obj *CodeLinterOption
		err field.ErrorList
	)
	BeforeEach(func() {
		obj = &CodeLinterOption{}
		err = field.ErrorList{}
	})
	JustBeforeEach(func() {
		err = obj.Validate(field.NewPath("test"))
	})
	Context("quality gate is not enabled", func() {
		BeforeEach(func() {
			obj.QualityGate = false
		})
		It("no error is returned", func() {
			Expect(err).To(BeEmpty())
		})
	})
	Context("quality gate is enabled", func() {
		BeforeEach(func() {
			obj.QualityGate = true
		})
		Context("when no quality gate rules", func() {
			It("no error is returned", func() {
				obj.QualityGateRules = nil
				Expect(err).To(BeEmpty())
			})
		})
		Context("when quality gate rules is not empty", func() {
			When("rule is valid", func() {
				BeforeEach(func() {
					obj.QualityGateRules = map[string]string{
						gateRuleIssuesCount: "10",
					}
				})
				It("no error is returned", func() {
					Expect(err).To(BeEmpty())
				})
			})
			When("rule is invalid", func() {
				BeforeEach(func() {
					obj.QualityGateRules = map[string]string{
						gateRuleIssuesCount: "10a",
					}
				})
				It("error is returned", func() {
					Expect(err).ShouldNot(BeEmpty())
				})
			})
		})
	})
})

var _ = Describe("Test.CodeLinterOption.ValidateQualityGate", func() {
	var opt *CodeLinterOption
	var ctx = context.Background()
	BeforeEach(func() {
		opt = &CodeLinterOption{}
	})
	Context("enable quality gate", func() {
		BeforeEach(func() {
			opt.QualityGate = true
		})
		When("quality gate rules is empty", func() {
			BeforeEach(func() {
				opt.QualityGateRules = nil
			})
			It("should ignore quality check and return nil", func() {
				opt.Issues = &v1alpha1.CodeLintIssues{Count: 11}
				errs := opt.ValidateQualityGate(ctx)
				Expect(errs).To(BeNil())
			})
		})
		When("quality gate rules is not empty", func() {
			BeforeEach(func() {
				opt.QualityGateRules = map[string]string{
					gateRuleIssuesCount: "10",
				}
			})
			When("issue count is greater than threshold value", func() {
				It("should return error", func() {
					opt.Issues = &v1alpha1.CodeLintIssues{Count: 11}
					errs := opt.ValidateQualityGate(ctx)
					Expect(errs).NotTo(BeNil())
				})
			})
			When("issue count is not greater than threshold value", func() {
				It("should return nil", func() {
					opt.Issues = &v1alpha1.CodeLintIssues{Count: 9}
					errs := opt.ValidateQualityGate(ctx)
					Expect(errs).To(BeNil())

					opt.Issues = &v1alpha1.CodeLintIssues{Count: 10}
					errs = opt.ValidateQualityGate(ctx)
					Expect(errs).To(BeNil())
				})
			})
		})
	})
	Context("disable quality gate", func() {
		JustBeforeEach(func() {
			opt.QualityGate = false
		})
		It("should return nil", func() {
			errs := opt.ValidateQualityGate(ctx)
			Expect(errs).To(BeNil())
		})
	})
})

var _ = Describe("Test.CodeLinterOption.WriteResult", func() {
	var obj *CodeLinterOption
	var writer io.Writer
	BeforeEach(func() {
		obj = &CodeLinterOption{
			CodeLintResult: v1alpha1.CodeLintResult{
				Result: "Successed",
				Issues: &v1alpha1.CodeLintIssues{Count: 10},
			},
		}
		writer = &strings.Builder{}
	})

	Context("write result with error", func() {
		When("is quality gate check failed error", func() {
			It("should get fail result", func() {
				obj.WriteResult(qualitygate.QualityGateCheckFailedErr, writer)
				Expect(writer.(*strings.Builder).String()).To(MatchJSON(`{"result":"Failed","issues.count":"10"}`))
			})
		})
		When("is known error", func() {
			It("should get empty result", func() {
				obj.WriteResult(errors.New("test error"), writer)
				Expect(writer.(*strings.Builder).String()).To(Equal(""))
			})
		})
	})
	Context("write result without error", func() {
		It("should get success result", func() {
			obj.WriteResult(nil, writer)
			Expect(writer.(*strings.Builder).String()).To(MatchJSON(`{"result":"Successed","issues.count":"10"}`))
		})
	})
})
