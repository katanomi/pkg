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

package watcher

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/emicklei/go-restful/v3"
	kclient "github.com/katanomi/pkg/client"
	"github.com/katanomi/pkg/plugin/route"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	ktesting "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	cw         *CertWatcher
	secretName = "test-secret"
	namespace  = "katanomi-system"
)

func TestMain(m *testing.M) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	logger := zap.NewRaw(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)).Sugar()

	os.Setenv(EnvCertSecret, secretName)
	os.Setenv(EnvPodNamespace, namespace)

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	ctx := context.TODO()
	clt := fake.NewClientBuilder().WithScheme(scheme).Build()
	ctx = kclient.WithDirectClient(ctx, clt)
	ctx = kclient.WithClient(ctx, clt)
	ctx = logging.WithLogger(ctx, logger)
	config := ctrl.GetConfigOrDie()
	container := restful.NewContainer()
	container.Add(route.NewDefaultService())

	ktesting.LoadKubeResources("testdata/secret.yaml", clt)

	cw = NewCertWatcher(ctx, config, container, "testdata/cert")

	os.Exit(m.Run())
}

func TestGetCAAndRemoveFile(t *testing.T) {
	g := NewGomegaWithT(t)

	ca := cw.GetCA()
	g.Expect(ca).NotTo(BeNil())

	err := cw.WaitCertFilesCreation()
	g.Expect(err).To(BeNil())

	// check tls.key file is delete
	t.Log("wait exec remove func...")
	time.Sleep(time.Second * 65)
	_, err = os.Stat(cw.keyPath)
	g.Expect(os.IsNotExist(err)).To(BeTrue())
}

func TestUpdateCertKey(t *testing.T) {
	g := NewGomegaWithT(t)

	secret := &corev1.Secret{}
	clt := kclient.DirectClient(cw.ctx)

	ns, err := podNamespace()
	g.Expect(err).To(BeNil())

	name, err := certSecretName()
	g.Expect(err).To(BeNil())

	key := types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}
	err = clt.Get(cw.ctx, key, secret)
	g.Expect(err).To(BeNil())

	time.AfterFunc(time.Second*10, shotdownServer)
	updateErr := cw.updateCertKey(secret)
	g.Expect(updateErr).To(Equal(http.ErrServerClosed))
}

func shotdownServer() {
	fmt.Println("--- shotdown server ---")
	cw.Server.Shutdown(cw.ctx)
}
