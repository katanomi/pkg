/*
Copyright 2023 The Katanomi Authors.

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

package workers

import (
	"context"
	"os"

	"github.com/go-resty/resty/v2"
	kconfig "github.com/katanomi/pkg/config"
	"github.com/katanomi/pkg/restclient"
	client2 "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	mockmgr "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"knative.dev/pkg/configmap/informer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type fakeRunner struct {
	result string
}

func (f *fakeRunner) Setup(ctx context.Context, kclient client.Client, restClient *resty.Client) error {
	return nil
}

func (f *fakeRunner) JobName() string {
	return "fake-runner"
}

func (f *fakeRunner) RunFunc(ctx context.Context) func() {
	return func() {
		f.result = "ok"
	}
}

var _ = Describe("Test.CronWorker", func() {
	var (
		ctx     context.Context
		manager *mockmgr.MockManager
		watcher *informer.InformedWatcher
		kclient client.Client
		cfgMgr  *kconfig.Manager
		logger  *zap.SugaredLogger

		fr     *fakeRunner
		worker *CronWorker

		cmName string
		ns     string

		stopChan chan struct{}
	)

	BeforeEach(func() {
		stopChan = make(chan struct{})
		ConfigWatcherFunc = func(cw *CronWorker) func(c *kconfig.Config) {
			return func(c *kconfig.Config) {
				for _, job := range cw.Runners {
					job.RunFunc(ctx)()
				}
				stopChan <- struct{}{}
			}
		}

		ctx = context.Background()
		mockCtl := gomock.NewController(GinkgoT())
		manager = mockmgr.NewMockManager(mockCtl)
		kclient = client2.NewMockClient(mockCtl)
		manager.EXPECT().GetClient().Return(kclient)
		manager.EXPECT().Add(gomock.AssignableToTypeOf(&CronWorker{})).Return(nil)

		logger = zap.NewNop().Sugar()

		ns = "default"
		cmName = "cm"
		oldCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: ns,
			},
		}

		client := k8sfake.NewSimpleClientset(oldCM)
		os.Setenv("SYSTEM_NAMESPACE", ns)
		watcher = informer.NewInformedWatcher(client, ns)
		cfgMgr = kconfig.NewManager(watcher, logger, cmName)

		ctx = kconfig.WithKatanomiConfigManager(ctx, cfgMgr)
		ctx = restclient.WithRESTClient(ctx, restyClient)
		fr = &fakeRunner{}

		worker = &CronWorker{
			Runners: []JobRunnable{
				fr,
			},
		}
	})

	Context("cron worker start", func() {

		It("run job successfully", func() {
			done := make(chan struct{})
			go func() {
				done <- struct{}{}
				stopCh := make(chan struct{})
				if err := watcher.Start(stopCh); err != nil {
					logger.Fatal("failed to start watcher", zap.Error(err))
				}
			}()
			<-done
			Expect(worker.Setup(ctx, manager, logger)).To(Succeed())
			watcher.OnChange(&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: ns,
				},
				Data: map[string]string{
					"xx": "yy",
				},
			})
			<-stopChan
			Expect(fr.result).To(Equal("ok"))
		})
	})
})
