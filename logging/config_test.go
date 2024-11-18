/*
Copyright 2024 The AlaudaDevops Authors.

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

package logging

import (
	"context"

	"github.com/AlaudaDevops/pkg/testing"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/logging"
)

var _ = Describe("LevelManager_Update", func() {
	var lm LevelManager
	var cm = &corev1.ConfigMap{}
	testing.MustLoadYaml("testdata/config.yaml", cm)

	BeforeEach(func() {
		ctx := logging.WithLogger(context.Background(), logger)
		lm = NewLevelManager(ctx, "test")
	})

	JustBeforeEach(func() {
		lm.Update()(cm)
	})

	Context("test for AtomicLevel", func() {
		var infoLevel zap.AtomicLevel
		var errorLevel zap.AtomicLevel
		BeforeEach(func() {
			infoLevel = lm.Get("test_info")
			errorLevel = lm.Get("test_error")
		})

		It("should get correct log level", func() {
			Expect(infoLevel.Level()).To(Equal(zap.InfoLevel))
			Expect(errorLevel.Level()).To(Equal(zap.ErrorLevel))
		})
	})

	Context("test for custom observer", func() {
		var gotData map[string]string
		BeforeEach(func() {
			lm.AddCustomObserver(func(d map[string]string) {
				gotData = d
			})
		})

		It("should get full data of the configmap", func() {
			expectCm := &corev1.ConfigMap{}
			testing.MustLoadYaml("testdata/config.yaml", expectCm)
			Expect(cmp.Diff(expectCm.Data, gotData)).To(BeEmpty())
		})
	})
})
