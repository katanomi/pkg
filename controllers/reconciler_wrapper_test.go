/*
Copyright 2023 The AlaudaDevops Authors.

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

package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	mockclient "github.com/AlaudaDevops/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilerWrapperReconcile(t *testing.T) {

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "test",
			Name:      "test",
		},
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
		},
	}

	cases := map[string]struct {
		r            reconcile.Reconciler
		requestFuncs []RequestFunc
		resultFuncs  []ResultFunc
		ctx          map[string]string
		result       reconcile.Result
		expect       func(*mockclient.MockClient)
		eval         func(g Gomega, err error)
	}{
		"nil reconciler": {
			r: nil,
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Times(0)
			},
			eval: func(g Gomega, err error) {
				g.Expect(err.Error()).To(ContainSubstring("reconciler should not be empty"))
			},
		},
		"empty request func": {
			r:            reconcile.Func(emptyReconciler),
			requestFuncs: []RequestFunc{emptyRequest},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Times(1)
			},
			result: reconcile.Result{},
		},
		"object not found": {
			r:            reconcile.Func(emptyReconciler),
			requestFuncs: []RequestFunc{emptyRequest},
			result:       reconcile.Result{},
			expect: func(c *mockclient.MockClient) {
				err := errors.NewNotFound(schema.GroupResource{}, "test")
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(err).Times(1)
			},
		},
		"request with error": {
			r:            reconcile.Func(emptyReconciler),
			requestFuncs: []RequestFunc{requestWithError},
			result:       reconcile.Result{},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
			eval: func(g Gomega, err error) {
				g.Expect(err.Error()).To(ContainSubstring("func return error"))
			},
		},
		"request with context": {
			r: reconcile.Func(emptyReconciler),
			ctx: map[string]string{
				"abc": "def",
			},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
			result: reconcile.Result{},
		},
		"multi request": {
			r: reconcile.Func(emptyReconciler),
			ctx: map[string]string{
				"abc": "def",
				"def": "abc",
			},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
			result: reconcile.Result{},
		},
		"empty result func": {
			r:           reconcile.Func(emptyReconciler),
			resultFuncs: []ResultFunc{emptyResult},
			result:      reconcile.Result{},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
		},
		"result with error": {
			r:           reconcile.Func(emptyReconciler),
			resultFuncs: []ResultFunc{resultWithError},
			result:      reconcile.Result{},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
			eval: func(g Gomega, err error) {
				g.Expect(err.Error()).To(ContainSubstring("func return error"))
			},
		},
		"multi result": {
			r:           reconcile.Func(emptyReconciler),
			resultFuncs: []ResultFunc{resultWithRequeue, resultWithRequeueAfter},
			result: reconcile.Result{
				Requeue:      true,
				RequeueAfter: 2 * time.Minute,
			},
			expect: func(c *mockclient.MockClient) {
				c.EXPECT().Get(gomock.Any(), request.NamespacedName, cm).Return(nil).Times(1)
			},
		},
	}

	for name, item := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			for key, value := range item.ctx {
				item.requestFuncs = append(item.requestFuncs, requestWithContext(key, value))
				item.resultFuncs = append(item.resultFuncs, resultHashContextValue(g, key, value))
			}

			ctrl := gomock.NewController(t)
			mockClient := mockclient.NewMockClient(ctrl)
			item.expect(mockClient)

			r := NewReconcilerWrapper(item.r, func(options *WrapperOptions) {
				options.RequestFuncs = item.requestFuncs
				options.ResultFuncs = item.resultFuncs
				options.Client = mockClient
				options.Object = cm
			})

			result, err := r.Reconcile(context.Background(), request)
			if item.eval != nil {
				item.eval(g, err)
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(result).To(Equal(item.result))
		})
	}
}

func emptyReconciler(_ context.Context, _ reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func emptyRequest(ctx context.Context, _ reconcile.Request) (context.Context, error) {
	return ctx, nil
}

func requestWithError(ctx context.Context, _ reconcile.Request) (context.Context, error) {
	return ctx, fmt.Errorf("request err")
}

func emptyResult(_ context.Context, _ reconcile.Request, result *reconcile.Result) error {
	return nil
}

func resultWithError(_ context.Context, _ reconcile.Request, result *reconcile.Result) error {
	return fmt.Errorf("result error")
}

func resultWithRequeue(_ context.Context, _ reconcile.Request, result *reconcile.Result) error {
	result.Requeue = true

	return nil
}

func resultWithRequeueAfter(_ context.Context, _ reconcile.Request, result *reconcile.Result) error {
	result.RequeueAfter = 2 * time.Minute

	return nil
}

func requestWithContext(key, value string) RequestFunc {
	return func(ctx context.Context, _ reconcile.Request) (context.Context, error) {
		return context.WithValue(ctx, key, value), nil
	}
}

func resultHashContextValue(g Gomega, key, value string) ResultFunc {
	return func(ctx context.Context, _ reconcile.Request, result *reconcile.Result) error {
		g.Expect(ctx.Value(key)).To(Equal(value))
		return nil
	}
}
