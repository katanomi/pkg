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

package io

import (
	"path"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWriteFile(t *testing.T) {
	g := NewGomegaWithT(t)

	dstFile := path.Join("./testdata", "dest", "test-write-file")
	g.Expect(WriteFile(dstFile, []byte("hello"), 0644)).Should(Succeed())
}

func TestFile(t *testing.T) {
	g := NewGomegaWithT(t)
	testFile := path.Join("./testdata", "src", "source-file")
	testDir := "./testdata"

	g.Expect(IsFile(testFile)).To(BeTrue())
	g.Expect(IsFile(testDir)).To(BeFalse())

	g.Expect(IsDir(testDir)).To(BeTrue())
	g.Expect(IsDir(testFile)).To(BeFalse())

	g.Expect(IsExist(testDir)).To(BeTrue())
	g.Expect(IsExist(testFile)).To(BeTrue())
	g.Expect(IsExist(path.Join("./testdata", "not-exist"))).To(BeFalse())
}
