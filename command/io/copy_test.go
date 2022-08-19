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
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCopy(t *testing.T) {
	g := NewGomegaWithT(t)

	cleanup := func() {
		err := os.RemoveAll("./testdata/dest")
		g.Expect(err).Should(Succeed())

		err = os.MkdirAll("./testdata/dest", 0777)
		g.Expect(err).Should(Succeed())
	}

	checkContent := func(file string, expectContent string) {
		fileContent, err := ioutil.ReadFile(file)
		g.Expect(err).Should(Succeed())
		g.Expect(bytes.Contains(fileContent, []byte(expectContent))).To(BeTrue())
	}

	checkDstFileContent := func(fileName string) {
		dstFile := path.Join("./testdata/dest", fileName)
		g.Expect(IsFile(dstFile)).To(BeTrue())
		checkContent(dstFile, "hello")
	}

	cleanup()
	defer cleanup()

	// copy dir to dir
	err := Copy("./testdata/src", "./testdata/dest")
	g.Expect(err).Should(Succeed())
	checkDstFileContent("source-file")
	checkDstFileContent("..source-file-link1")
	checkDstFileContent("..source-file-link2")
	checkDstFileContent("..source-file-link3")

	// copy file to file
	dstFile4 := path.Join("./testdata/dest", "dst-file-1")
	err = Copy("./testdata/src/source-file", dstFile4)
	g.Expect(err).Should(Succeed())
	checkDstFileContent("dst-file-1")

	dstFile5 := path.Join("./testdata/dest", "..dst-file-2")
	err = Copy("./testdata/src/source-file", dstFile5, FileFilterOption(func(srcFile, dstFile string) bool {
		if strings.HasPrefix(path.Base(dstFile), "..") {
			return true
		}
		return false
	}))
	g.Expect(err).Should(Succeed())
	g.Expect(IsExist(dstFile5)).To(BeFalse())

	dstFile6 := path.Join("./testdata/dest", "..dst-file-3")
	err = Copy("./testdata/src/source-file", dstFile6, FileFilterOption(func(srcFile, dstFile string) bool {
		if strings.HasPrefix(path.Base(dstFile), "xx") {
			return true
		}
		return false
	}))
	g.Expect(err).Should(Succeed())
	checkDstFileContent("..dst-file-3")
}
