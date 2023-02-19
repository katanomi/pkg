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

package config

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kclient "github.com/katanomi/pkg/client"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/system"
	_ "knative.dev/pkg/system/testing"
	kfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("NewManger and GetConfig and GetFeatureFlags", func() {

	var (
		configmap  *corev1.ConfigMap
		logger     *zap.Logger
		manager    *Manager
		watcher    *informer.InformedWatcher
		ns, cmName string
		stopCh     chan struct{}

		data map[string]string
	)
	BeforeEach(func() {
		ns = "default"
		cmName = "cm"
		configmap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: ns,
			},
		}
		logger, _ = zap.NewDevelopment()

		client := fake.NewSimpleClientset(configmap)
		watcher = informer.NewInformedWatcher(client, ns)

		manager = NewManager(watcher, logger.Sugar(), cmName)

		data = map[string]string{}
		stopCh = make(chan struct{})
		if err := watcher.Start(stopCh); err != nil {
			logger.Fatal("failed to start watcher", zap.Error(err))
		}
	})

	JustAfterEach(func() {
		close(stopCh)
	})

	When("create a manager", func() {
		When("get config before configmap change", func() {
			It("return empty data", func() {
				Expect(manager.GetConfig().Data).To(HaveLen(0))
			})
		})

		When("get config before configmap change", func() {
			It("return default feature flags", func() {
				Expect(manager.GetFeatureFlag(VersionEnabledFeatureKey)).To(Equal(DefaultVersionEnabled))
				Expect(manager.GetFeatureFlag(InitializeAllowLocalRequestsFeatureKey)).To(Equal(DefaultInitializeAllowLocalRequests))
				Expect(manager.GetFeatureFlag(PrunerDelayAfterCompletedFeatureKey)).To(Equal(DefaultPrunerDelayAfterCompleted))
				Expect(manager.GetFeatureFlag(PrunerKeepFeatureKey)).To(Equal(DefaultPrunerKeep))
			})
		})

		When("get config after update configmap ", func() {
			It("should return updated data", func() {
				By("change the manger watched configmap")
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
				Expect(manager.GetConfig().Data).To(Equal(data), "should return updated config")

				By("change the manger watched configmap again")
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
				Expect(manager.GetConfig().Data).To(Equal(data), "should return updated config again")
			})
		})

		When("get config after update configmap ", func() {
			It("should return updated data", func() {
				By("change the manger watched configmap")

				data := map[string]string{
					VersionEnabledFeatureKey: "true",
					PrunerKeepFeatureKey:     "6000",
					"other.test":             "1",
				}
				watcher.OnChange(&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cmName,
						Namespace: ns,
					},
					Data: data,
				})
				Expect(manager.GetFeatureFlag(VersionEnabledFeatureKey)).To(Equal(FeatureValue("true")))
				Expect(manager.GetFeatureFlag(PrunerKeepFeatureKey)).To(Equal(FeatureValue("6000")))
				Expect(manager.GetFeatureFlag(InitializeAllowLocalRequestsFeatureKey)).To(Equal(DefaultInitializeAllowLocalRequests))
				Expect(manager.GetFeatureFlag(PrunerDelayAfterCompletedFeatureKey)).To(Equal(DefaultPrunerDelayAfterCompleted))

				By("change the manger watched configmap again")

				data = map[string]string{
					VersionEnabledFeatureKey: "false",
					PrunerKeepFeatureKey:     "1000",
					"other.test":             "1",
				}

				watcher.OnChange(&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      cmName,
						Namespace: ns,
					},
					Data: data,
				})
				Expect(manager.GetFeatureFlag(VersionEnabledFeatureKey)).To(Equal(FeatureValue("false")))
				Expect(manager.GetFeatureFlag(PrunerKeepFeatureKey)).To(Equal(FeatureValue("1000")))
				Expect(manager.GetFeatureFlag(InitializeAllowLocalRequestsFeatureKey)).To(Equal(DefaultInitializeAllowLocalRequests))
				Expect(manager.GetFeatureFlag(PrunerDelayAfterCompletedFeatureKey)).To(Equal(DefaultPrunerDelayAfterCompleted))
				Expect(manager.GetFeatureFlag("other.test")).To(Equal(FeatureValue("1")))

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

func TestIsSameConfigMap(t *testing.T) {
	t.Run("objectmeta match, return true", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()

		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		sameConfig := corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cm", Namespace: system.Namespace()}}
		g.Expect(manager.isSameConfigMap(&sameConfig)).To(Equal(true))
	})

	t.Run("when objectmeta not match, return false", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()

		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		sameConfig := corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: system.Namespace()}}
		g.Expect(manager.isSameConfigMap(&sameConfig)).To(Equal(false))
	})

	t.Run("when objectmeta not match, return false", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()

		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		manager = nil
		sameConfig := corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: system.Namespace()}}
		g.Expect(manager.isSameConfigMap(&sameConfig)).To(Equal(false), "when manager is nil, return false")
	})
}

func Test_GetFeatureFlagByClient(t *testing.T) {
	t.Run("manager is empty", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var manager *Manager
		got := manager.GetFeatureFlagByClient(context.TODO(), VersionEnabledFeatureKey)
		g.Expect(got).To(Equal(FeatureValue("false")))
	})

	t.Run("context client is empty", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()

		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		got := manager.GetFeatureFlagByClient(context.TODO(), VersionEnabledFeatureKey)
		g.Expect(got).To(Equal(FeatureValue("false")))
	})

	t.Run("configmap object not found", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()
		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		ctx := context.TODO()
		clt := kfake.NewClientBuilder().Build()
		ctx = kclient.WithClient(ctx, clt)

		got := manager.GetFeatureFlagByClient(ctx, VersionEnabledFeatureKey)
		g.Expect(got).To(Equal(FeatureValue("false")))
	})

	t.Run("get feature flag success.", func(t *testing.T) {
		g := NewGomegaWithT(t)
		client := fake.NewSimpleClientset()

		watcher := informer.NewInformedWatcher(client, "cm")
		manager := NewManager(watcher, nil, "cm")

		cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cm", Namespace: system.Namespace()}, Data: map[string]string{"flag": "true"}}
		ctx := context.TODO()
		clt := kfake.NewClientBuilder().WithObjects(cm).Build()
		ctx = kclient.WithClient(ctx, clt)

		got := manager.GetFeatureFlagByClient(ctx, "flag")
		g.Expect(got).To(Equal(FeatureValue("true")))
	})
}
