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

package framework

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// NewServicePortForward construct a new ServicePortForward object
func NewServicePortForward(ns string, serviceName string) *ServicePortForward {
	return &ServicePortForward{
		Namespace:   ns,
		ServiceName: serviceName,
	}
}

// ServicePortForward describe the configuration of port forwarding
type ServicePortForward struct {
	Namespace         string
	ServiceName       string
	ServicePortFilter func(ports []v1.ServicePort) int32

	stopChan chan struct{}
}

// Close will stop the forwarder
func (s *ServicePortForward) Close() {
	if s.stopChan != nil {
		close(s.stopChan)
	}
}

// filterServicePort filter the service port
func (s *ServicePortForward) filterServicePort(ports []v1.ServicePort) int32 {
	if s.ServicePortFilter != nil {
		return s.ServicePortFilter(ports)
	}

	if len(ports) > 0 {
		return int32(ports[0].TargetPort.IntValue())
	}

	return 0
}

// Forward setting a forwarder
func (s *ServicePortForward) Forward(testCtx *TestContext) (restyClient *resty.Client, err error) {
	ctx := testCtx.Context
	clt := testCtx.Client
	logger := testCtx.With("Namespace", s.Namespace, "ServiceName", s.ServiceName)

	service := &v1.Service{}
	serviceKey := types.NamespacedName{Namespace: s.Namespace, Name: s.ServiceName}
	if err = clt.Get(ctx, serviceKey, service); err != nil {
		logger.Errorw("get service error", "err", err)
		return nil, err
	}

	servicePort := s.filterServicePort(service.Spec.Ports)

	podList := &v1.PodList{}
	podLabel := client.MatchingLabels(service.Spec.Selector)
	if err = clt.List(ctx, podList, podLabel, client.InNamespace(s.Namespace)); err != nil {
		logger.Errorw("get pod list error", "err", err)
		return nil, err
	}

	if len(podList.Items) == 0 {
		err = errors.New("no pods to choose from")
		logger.Errorw("get service pod error", "err", err)
		return nil, err
	}

	pod := podList.Items[0]
	return s.forward(testCtx, pod.GetName(), servicePort)
}

func (s *ServicePortForward) forward(testCtx *TestContext, podName string, forwardPort int32) (restyClient *resty.Client, err error) {
	var (
		config = testCtx.Config
		logger = testCtx.With("Namespace", s.Namespace, "podName", podName)
	)
	clientSet := clientset.NewForConfigOrDie(config)
	req := clientSet.CoreV1().RESTClient().Post().Namespace(s.Namespace).
		Resource("pods").Name(podName).SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return nil, err
	}

	s.stopChan = make(chan struct{})
	readyChan := make(chan struct{})

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req.URL())
	ports := []string{fmt.Sprintf(":%d", forwardPort)}
	pf, err := portforward.New(dialer, ports, s.stopChan, readyChan, os.Stdout, os.Stderr)
	if err != nil {
		logger.Errorw("forward port error", "err", err)
		return nil, err
	}

	go func() {
		if forwardErr := pf.ForwardPorts(); forwardErr != nil {
			logger.Errorw("forwardPorts returned with error", "err", forwardErr)
		}
	}()

	<-pf.Ready

	// ignore the error because portforward should be ready now
	list, _ := pf.GetPorts()
	if len(list) == 0 || list[0].Local == 0 {
		err = errors.New("get local port error")
		logger.Errorw("get local port error", "portList", list)
		return nil, err
	}

	restyClient = resty.New().SetHostURL(fmt.Sprintf("http://localhost:%d", list[0].Local))
	return restyClient, err
}
