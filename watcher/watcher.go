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

// Package watcher stores interface DefaultingWatcherWithOnChange
package watcher

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/emicklei/go-restful/v3"
	kclient "github.com/katanomi/pkg/client"
	kscheme "github.com/katanomi/pkg/scheme"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	v1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/logging"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"yunion.io/x/pkg/util/wait"
)

const (
	EnvPodNamespace = "SYSTEM_NAMESPACE"
	EnvCertSecret   = "CERT_SECRET"
)

type CertWatcher struct {
	sync.Mutex
	ctx        context.Context
	config     *rest.Config
	secretName string

	caPath  string
	crtPath string
	keyPath string

	Logger *zap.SugaredLogger

	container *restful.Container
	Server    *http.Server

	removeTimer *time.Timer
}

func NewCertWatcher(ctx context.Context, config *rest.Config, container *restful.Container, certPath string) *CertWatcher {
	logger := logging.FromContext(ctx)
	secretName, err := certSecretName()
	if err != nil {
		logger.Fatalf("get secret error", "error", err)
	}
	secret := &corev1.Secret{}
	clt := kclient.DirectClient(ctx)
	if clt == nil {
		newClt, err := ctrlclient.New(config, ctrlclient.Options{Scheme: kscheme.Scheme(ctx)})
		if err != nil {
			logger.Fatalf("create new k8s client error", "error", err)
		}
		clt = newClt
	}
	podNS, err := podNamespace()
	if err != nil {
		logger.Fatalf("get pod namespace error", "error", err)
	}
	if err := clt.Get(ctx, types.NamespacedName{Namespace: podNS, Name: secretName}, secret); err != nil {
		logger.Fatalf("NewCertWatcher get cert secret error", "error", err)
	}

	_, err = os.Stat(certPath)
	if err != nil && !os.IsNotExist(err) {
		logger.Fatalf("stat cert path error", "error", err)
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(certPath, 0700)
		if err != nil {
			logger.Fatalf("create cert path error", "error", err)
		}
	}

	cw := &CertWatcher{
		ctx:        ctx,
		config:     config,
		secretName: secretName,
		container:  container,
		caPath:     filepath.Join(certPath, "ca.crt"),
		crtPath:    filepath.Join(certPath, "tls.crt"),
		keyPath:    filepath.Join(certPath, "tls.key"),
		Logger:     logger,
	}

	if err := cw.writeCert(secret); err != nil {
		logger.Fatalf("NewCertWatcher set secret error", "error", err)
	}

	return cw
}

func (cw *CertWatcher) Start() (err error) {
	cw.Logger.Info("start to listen cert secret change")
	c, err := kubernetes.NewForConfig(cw.config)
	if err != nil {
		return
	}

	podNS, err := podNamespace()
	if err != nil {
		cw.Logger.Errorw("get pod namespace error", "error", err)
		return err
	}
	informer := v1.NewFilteredSecretInformer(c, podNS, time.Hour*12,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, func(options *metav1.ListOptions) {
			options.FieldSelector = fmt.Sprintf("metadata.name=%s", cw.secretName)
		})
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			err = cw.updateCertKey(obj.(*corev1.Secret))
			if err != nil {
				cw.Logger.Errorw("load secret error", "error", err)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			err = cw.updateCertKey(newObj.(*corev1.Secret))
			if err != nil {
				cw.Logger.Errorw("load secret error", "error", err)
			}
		},
		DeleteFunc: nil,
	})
	go informer.Run(cw.ctx.Done())
	cw.Logger.Infof("CertWatcher is started.")
	return
}

func (cw *CertWatcher) updateCertKey(secret *corev1.Secret) (err error) {
	cw.Logger.Infof("set cert with secret")
	if err := cw.writeCert(secret); err != nil {
		cw.Logger.Errorw("update cert key with secret error", "error", err)
		return err
	}

	if err = cw.RestartServer(); err != nil {
		cw.Logger.Errorw("start https server error", "error", err)
		return err
	}

	return nil
}

func (cw *CertWatcher) StartServer() (err error) {
	cw.Logger.Info("start https server")
	port := 8443
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: cw.container,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		},
	}
	cw.Server = srv
	return srv.ListenAndServeTLS(cw.crtPath, cw.keyPath)
}

func (cw *CertWatcher) RestartServer() (err error) {
	cw.Logger.Warn("Restart https server")
	if cw.Server == nil {
		return cw.StartServer()
	}

	cw.Logger.Info("RestartServer shutdown https server")
	if err := cw.Server.Shutdown(cw.ctx); err != nil {
		cw.Logger.Errorw("shotdown server error", "error", err)
		return err
	}

	cw.Logger.Info("RestartServer start server again")
	return cw.StartServer()
}

func (cw *CertWatcher) writeCert(secret *corev1.Secret) error {
	cw.Logger.Debugw("set cert with secret", "certData", secret.Data)
	if string(secret.Data["tls.key"]) == "" || string(secret.Data["ca.crt"]) == "" || string(secret.Data["tls.crt"]) == "" {
		cw.Logger.Errorw("get cert data from secret error: cert is empty", "secret", secret)
		return fmt.Errorf("cert is empty")
	}

	if err := os.WriteFile(cw.keyPath, secret.Data["tls.key"], 0600); err != nil {
		cw.Logger.Errorw("write tls.key error", "error", err)
		return err
	}
	if err := os.WriteFile(cw.caPath, secret.Data["ca.crt"], 0600); err != nil {
		cw.Logger.Errorw("write ca.crt error", "error", err)
		return err
	}
	if err := os.WriteFile(cw.crtPath, secret.Data["tls.crt"], 0600); err != nil {
		cw.Logger.Errorw("write tls.crt error", "error", err)
		return err
	}

	cw.Logger.Infof("Update cert files success.")
	cw.Lock()
	defer cw.Unlock()
	if cw.removeTimer != nil {
		cw.removeTimer.Stop()
	}
	cw.removeTimer = time.AfterFunc(time.Second*60, cw.removeKeyFile)

	return nil
}

func (cw *CertWatcher) removeKeyFile() {
	cw.Lock()
	defer cw.Unlock()
	err := os.Remove(cw.keyPath)
	if err != nil && !os.IsNotExist(err) {
		cw.Logger.Errorw("Remove key file failed", "error", err)
	}
	cw.Logger.Infof("The webhook key file is removed.")
	cw.removeTimer = nil
}

func (cw *CertWatcher) WaitCertFilesCreation() error {
	return wait.Poll(time.Millisecond*100, time.Second*30, func() (done bool, err error) {
		_, err = os.Stat(cw.keyPath)
		if err != nil {
			cw.Logger.Errorw("stat tls.key error", "error", err)
			return
		}
		_, err = os.Stat(cw.caPath)
		if err != nil {
			cw.Logger.Errorw("stat tls.key error", "error", err)
			return
		}
		_, err = os.Stat(cw.crtPath)
		if err != nil {
			cw.Logger.Errorw("stat tls.key error", "error", err)
			return
		}
		cw.Logger.Infof("Webhook cert files loaded.")
		return true, nil
	})
}

func (cw *CertWatcher) SetContainer(container *restful.Container) {
	cw.container = container
}

func (cw *CertWatcher) GetCA() []byte {
	file, err := ioutil.ReadFile(cw.caPath)
	if err != nil {
		cw.Logger.Fatalw("read ca file error", "error", err)
		return nil
	}
	return file
}

func podNamespace() (string, error) {
	ns := os.Getenv(EnvPodNamespace)
	if ns == "" {
		return "", fmt.Errorf("get pod namespace from env error")
	}
	return ns, nil
}

func certSecretName() (string, error) {
	name := os.Getenv(EnvCertSecret)
	if name == "" {
		return "", fmt.Errorf("get secret name from env error")
	}
	return name, nil
}
