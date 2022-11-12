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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/katanomi/pkg/encoding"
	testing2 "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
)

func TestVulnScanResult_IsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)
	empty := VulnScanResult{}
	g.Expect(empty.IsEmpty()).To(BeTrue())

	notEmpty := VulnScanResult{
		Result: "Successed",
	}
	g.Expect(notEmpty.IsEmpty()).To(BeFalse())
}

func TestVulnScanResult_decode(t *testing.T) {
	g := NewGomegaWithT(t)
	result := VulnScanResult{}

	m := make(map[string]string)
	testing2.MustLoadYaml("./testdata/VulnScanResult.decode.yaml", &m)
	g.Expect(encoding.NewJsonPath().Decode(&result, m)).Should(Succeed())

	goldenData := VulnScanResult{}
	testing2.MustLoadYaml("./testdata/VulnScanResult.decode.golden.yaml", &goldenData)

	diff := cmp.Diff(goldenData, result)
	g.Expect(diff).To(gomega.BeEmpty())
}

func TestVulnScanResult_encode(t *testing.T) {
	g := NewGomegaWithT(t)
	result := VulnScanResult{}
	testing2.MustLoadYaml("./testdata/VulnScanResult.decode.golden.yaml", &result)
	encodeData := encoding.NewJsonPath().Encode(&result)

	goldenData := make(map[string]string)
	testing2.MustLoadYaml("./testdata/VulnScanResult.decode.yaml", &goldenData)
	diff := cmp.Diff(encodeData, goldenData)
	g.Expect(diff).To(gomega.BeEmpty())
}
