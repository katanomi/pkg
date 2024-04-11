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

package client

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestPluginClientContext(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	ctx := context.TODO()

	clt, ok := PluginClientFrom(ctx)
	g.Expect(ok).To(BeFalse())
	g.Expect(clt).To(BeNil())

	pluginClient := &PluginClient{}
	ctx = WithPluginClient(ctx, pluginClient)

	g.Expect(PluginClientValue((ctx))).To(Equal(pluginClient))
	clt, ok = PluginClientFrom(ctx)
	g.Expect(ok).To(BeTrue())
	g.Expect(clt).To(Equal(pluginClient))
}

func TestGitPluginClientContext(t *testing.T) {
	g := NewGomegaWithT(t)
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	ctx := context.TODO()

	clt, ok := GitPluginClientFrom(ctx)
	g.Expect(ok).To(BeFalse())
	g.Expect(clt).To(BeNil())

	gitPluginClient := &GitPluginClient{}
	ctx = WithGitPluginClient(ctx, gitPluginClient)

	g.Expect(GitPluginClientValue((ctx))).To(Equal(gitPluginClient))
	clt, ok = GitPluginClientFrom(ctx)
	g.Expect(ok).To(BeTrue())
	g.Expect(clt).To(Equal(gitPluginClient))
}
