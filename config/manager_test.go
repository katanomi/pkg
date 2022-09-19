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

package config

import (
	"context"
	"os"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"knative.dev/pkg/configmap/informer"

	. "github.com/onsi/gomega"
	_ "knative.dev/pkg/system/testing"
)

var _ = Describe("NewManger and GetConfig", func() {

	var (
		oldCM      *corev1.ConfigMap
		logger     *zap.Logger
		manager    *Manager
		watcher    *informer.InformedWatcher
		ns, cmName string

		data map[string]string
	)
	BeforeEach(func() {
		ns = "default"
		cmName = "cm"
		oldCM = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: ns,
			},
		}
		logger, _ = zap.NewDevelopment()

		client := fake.NewSimpleClientset(oldCM)
		watcher := informer.NewInformedWatcher(client, ns)

		manager = NewManager(watcher, logger.Sugar(), cmName)

		stopCh := make(chan struct{})
		defer close(stopCh)
		if err := watcher.Start(stopCh); err != nil {
			logger.Fatal("failed to start watcher", zap.Error(err))
		}

	})

	Context("When create a manager", func() {

		Describe("get config before configmap change", func() {
			It("return empty data", func() {
				Expect(manager.GetConfig().Data).To(Equal(data))
			})
		})

		Describe("get config after update configmap ", func() {
			By("change the manger watched configmap", func() {
				data = map[string]string{
					"first": "first",
				}
				watcher.OnChange(&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cmName,
						Namespace: ns,
					},
					Data: data,
				})
			})
			It("return updated config", func() {
				Expect(manager.GetConfig().Data).To(Equal(data))
			})

			By("change the manger watched configmap again", func() {
				data = map[string]string{
					"second": "second",
				}
				watcher.OnChange(&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cmName,
						Namespace: ns,
					},
					Data: data,
				})

			})
			It("return updated config again", func() {
				Expect(manager.GetConfig().Data).To(Equal(data))
			})
		})
	})
})

func TestConfigName(t *testing.T) {
	g := NewGomegaWithT(t)
	configName := Name()
	g.Expect(configName).To(Equal(defaultConfig))

	newName := "new"
	os.Setenv(configNameEnv, newName)
	g.Expect(newName).To(Equal(newName))
}

func TestKCMContext(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := context.TODO()
	client := fake.NewSimpleClientset()

	watcher := informer.NewInformedWatcher(client, "cm")
	manager := NewManager(watcher, nil, "cm")

	ctx = WithKatanomiConfigManager(ctx, manager)
	g.Expect(KatanomiConfigManager(ctx)).To(Equal(manager))
}
