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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ktesting "github.com/katanomi/pkg/testing"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func TestObjectToRecord(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	cm := &corev1.ConfigMap{}
	ktesting.MustLoadYaml("testdata/objectToRecord.cm.yaml", cm)
	wantRecord := Record{}
	ktesting.MustLoadYaml("testdata/objectToRecord.cm.golden.yaml", &wantRecord)
	gotRecord := ObjectToRecord(cm)
	g.Expect(cmp.Diff(gotRecord, wantRecord)).To(gomega.BeEmpty())
}
