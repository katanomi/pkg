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

package fake

import (
	"context"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/testing/mock/testing/fake"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandlerSetup(t *testing.T) {
	g := NewGomegaWithT(t)

	var (
		ctx    = context.Background()
		logger = zap.S()
	)

	h := PolicyHandler{}
	err := h.Setup(ctx, restful.Add, logger)
	g.Expect(err).NotTo(BeNil())
	g.Expect(restful.DefaultContainer.RegisteredWebServices()).To(HaveLen(0))

	ctrl := gomock.NewController(t)
	store := fake.NewMockStore(ctrl)
	store.EXPECT().Setup(ctx).Return(nil)
	h.Store = store

	err = h.Setup(ctx, restful.Add, logger)
	g.Expect(err).To(BeNil())
	g.Expect(restful.DefaultContainer.RegisteredWebServices()).To(HaveLen(1))
}
