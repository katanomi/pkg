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

package route

import (
	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test.Storage.Route.Service", func() {
	var (
		svcs []*restful.WebService
	)
	Context("NewService with file store plugin", func() {
		It("returns services with file store routes", func() {
			svcs, _ = NewServices(&fakeFileStorePlugin{})
			Expect(svcs).To(HaveLen(1))
			Expect(svcs[0].RootPath()).To(Equal("/storage/fake-filestore/file-store/v1alpha1"))
		})
	})

	Context("NewService with auth plugin", func() {
		It("returns services with core routes", func() {
			svcs, _ = NewServices(&fakeCorePlugin{})
			Expect(svcs).To(HaveLen(1))
			Expect(svcs[0].RootPath()).To(Equal("/storage/fake-core/core/v1alpha1"))
		})
	})
})
