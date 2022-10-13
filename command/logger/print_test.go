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

package logger

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func TestResultErrors_Error(t *testing.T) {
	g := NewGomegaWithT(t)
	var out []string
	ctx := genTestLogContext(&out)

	errs := utilerrors.NewAggregate([]error{
		fmt.Errorf("test error1"),
		fmt.Errorf("test error2"),
	})
	err := ResultErrors(ctx, errs, "success", "return error")
	g.Expect(err).ShouldNot(Succeed())
	g.Expect(err.Error()).To(Equal("2 return error"))
	g.Expect(out).To(Equal([]string{
		"==> 🛑  test error1",
		"==> 🛑  test error2",
	}))
}

func TestResultErrors_Success(t *testing.T) {
	g := NewGomegaWithT(t)
	var out []string
	ctx := genTestLogContext(&out)

	err := ResultErrors(ctx, nil, "success", "return error")
	g.Expect(err).Should(Succeed())
	g.Expect(out).To(HaveLen(1))
	g.Expect(out[0]).To(Equal("==> ✅  success"))
}

func genTestLogContext(output *[]string) context.Context {
	logger := zap.NewExample(zap.Hooks(func(entry zapcore.Entry) error {
		*output = append(*output, entry.Message)
		return nil
	})).Sugar()
	return WithLogger(context.Background(), logger)
}
