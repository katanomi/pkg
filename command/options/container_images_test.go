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

package options

import (
	"io/ioutil"
	"os"

	"github.com/google/go-cmp/cmp"
	artifacts "github.com/katanomi/pkg/apis/artifacts/v1alpha1"
	pkgtest "github.com/katanomi/pkg/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.ContainerImagesOption.Validate", func() {
	var (
		imageOption *ContainerImagesOption
		rootPath    = field.NewPath("root")
	)

	Context("validate valid image", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{
				ContainerImages: []string{
					"docker.io/centos:latest",
					"docker.io/centos",
					"docker.io/centos:",
					"docker.io/centos@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"docker.io/centos:latest@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"127.0.0.1:8080/centos@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"127.0.0.1/centos@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
				},
			}
		})
		It("should validate success", func() {
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(0))
		})
	})

	Context("validate invalid image", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{
				ContainerImages: []string{
					"docker.io/centos: test",
					"docker.io/centos@sha256:744c8b3d4c8f5b30a1a7",
					"docker.io/centos:latest@sha234:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"127.0.0.1:8080/centos@744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"127.0.0.1/centos:中文tag",
				},
			}
		})

		It("should return same err", func() {
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(len(imageOption.ContainerImages)))
		})
	})

	Context("validate required image", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{}
		})

		It("should return error", func() {
			imageOption.SetValueRequired(true)
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(1))
		})
	})

	Context("validate required tag", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{
				ContainerImages: []string{
					"docker.io/centos",
					"docker.io/centos:",
					"127.0.0.1:8080/centos@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
				},
			}
			imageOption.SetTagRequired(true)
		})

		It("should return err", func() {
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(len(imageOption.ContainerImages)))
		})
	})

	Context("validate without digest", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{
				ContainerImages: []string{
					"docker.io/centos:test:1",
					"127.0.0.1:8080/centos@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
				},
			}
			imageOption.SetWithoutDigest(true)
		})

		It("should return err", func() {
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(len(imageOption.ContainerImages)))
		})
	})

	Context("validate without digest and tag", func() {
		JustBeforeEach(func() {
			imageOption = &ContainerImagesOption{
				ContainerImages: []string{
					"127.0.0.1:8080/centos:v1@sha256:744c8b3d4c8f5b30a1a78c5e3893c4d3f793919d1e14bcaee61028931e9f9929",
					"127.0.0.1:8080/centos",
				},
			}
			imageOption.SetWithoutDigest(true)
			imageOption.SetTagRequired(true)
			imageOption.SetValueRequired(true)
		})

		It("should return err", func() {
			errs := imageOption.Validate(rootPath)
			Expect(len(errs)).To(Equal(len(imageOption.ContainerImages)))
		})
	})
})

var _ = Describe("ContainerImagesOption.WriteResult", func() {
	var (
		opts         *ContainerImagesOption
		artifactList []artifacts.URI
		resultFile   *os.File
		err          error
	)
	BeforeEach(func() {
		resultFile, err = ioutil.TempFile("/tmp/", "test-ContainerImagesOption.WriteResult-*")
		Expect(err).To(BeNil())
		opts = &ContainerImagesOption{
			ResultPath: resultFile.Name(),
		}
		DeferCleanup(func() {
			Expect(os.Remove(resultFile.Name())).To(Succeed())
		})
	})
	JustBeforeEach(func() {
		err = opts.WriteResult(artifactList)
	})
	JustAfterEach(func() {
	})
	When("list is full with data", func() {
		BeforeEach(func() {
			artifactList = []artifacts.URI{
				{
					Host:      "index.docker.io",
					Path:      "katanomi/controller",
					Digest:    "88008bc80503efb3d6c0a8c76fbda9e89067fc57c400c89901519984fc80ad93",
					Algorithm: artifacts.SHA256,
					Tag:       "v1.2.3",
				},
				{
					Host: "index.docker.io",
					Path: "katanomi/controller",
					Tag:  "latest",
				},
				{
					Host:      "docker.io",
					Path:      "katanomi/another",
					Digest:    "88008bc80503efb3d6c0a8c76fbda9e89067fc57c400c89901519984fc80ad93",
					Algorithm: artifacts.SHA256,
				},
			}
		})
		It("should write a json array with uri strings ", func() {
			result := []string{}
			Expect(pkgtest.LoadJSON(resultFile.Name(), &result)).To(Succeed())
			expected := []string{}
			Expect(pkgtest.LoadJSON("testdata/container-images.WriteResult.full.golden.json", &expected)).To(Succeed())
			Expect(cmp.Diff(result, expected)).To(BeEmpty())
		})
	})
	When("list is nil", func() {
		BeforeEach(func() {
			artifactList = nil
		})
		It("should write an empty json array", func() {
			result := []string{}
			Expect(pkgtest.LoadJSON(resultFile.Name(), &result)).To(Succeed())
			expected := []string{}
			Expect(pkgtest.LoadJSON("testdata/container-images.WriteResult.nil.golden.json", &expected)).To(Succeed())
			Expect(cmp.Diff(result, expected)).To(BeEmpty())
		})
	})

})
