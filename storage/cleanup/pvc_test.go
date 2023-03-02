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

package cleanup

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kclient "github.com/katanomi/pkg/client"
	pkgTesting "github.com/katanomi/pkg/testing"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	k8scheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("cleanupPVC clean pvc which mounted by taskrun", func() {

	Context("buildrun has clone task and clone task has pvc mounted", func() {
		scheme := runtime.NewScheme()

		pvc := &corev1.PersistentVolumeClaim{}
		cachePVC := &corev1.PersistentVolumeClaim{}

		taskRun := &tekton.TaskRun{}
		pipelineRun := &tekton.PipelineRun{}

		BeforeEach(func() {
			_ = k8scheme.AddToScheme(scheme)
			_ = tekton.AddToScheme(scheme)

			err := pkgTesting.LoadYAML("../testdata/cleanup/pvc.yaml", pvc)
			Expect(err).To(BeNil())
			finalizers := pvc.GetFinalizers()
			Expect(finalizers).To(ContainElement(PVCFinalizersForPod))

			err = pkgTesting.LoadYAML("../testdata/cleanup/cache-pvc.yaml", cachePVC)
			Expect(err).To(BeNil())

			err = pkgTesting.LoadYAML("../testdata/cleanup/taskrun.yaml", taskRun)
			Expect(err).To(BeNil())
			err = pkgTesting.LoadYAML("../testdata/cleanup/pipelinerun.yaml", pipelineRun)
			Expect(err).To(BeNil())

		})

		It("clone related PVC will be deleted", func() {

			fakeClient := fakeclient.NewClientBuilder().WithScheme(scheme).WithObjects(cachePVC, pvc, taskRun, pipelineRun).Build()

			object := &corev1.PersistentVolumeClaim{}
			pvcKey := client.ObjectKey{
				Namespace: pvc.Namespace,
				Name:      pvc.Name,
			}
			cachePVCKey := client.ObjectKey{
				Namespace: cachePVC.Namespace,
				Name:      cachePVC.Name,
			}

			// reconcile should be ok
			ctx := kclient.WithClient(context.Background(), fakeClient)
			taskLables := client.MatchingLabels{"pipelineruns.tekton.dev/name": pipelineRun.Name}
			succeeded, fails, err := CleanTaskRunsPVC(ctx, taskLables,
				func(ctx context.Context, workspace tekton.WorkspaceBinding) bool {
					if workspace.PersistentVolumeClaim == nil {
						return true
					}
					if IsVolumeCreatedManaualy(pipelineRun.Spec.Workspaces, workspace) {
						return true
					}
					return false
				})
			Expect(err).To(BeNil())
			Expect(len(succeeded)).Should(Equal(1))
			Expect(len(fails)).Should(Equal(0))

			// pvc should be deleted
			err = fakeClient.Get(ctx, pvcKey, object)
			if err != nil {
				Expect(errors.IsNotFound(err)).To(BeTrue(), err.Error())
			} else {
				// pvc is not really removed, just removed finalizers
				// Ref: https://github.com/katanomi/builds/blob/b43dd50c07dcd83c008db51a135aa5a629792f45/pkg/controllers/builds/cleanup.go#L210-L212
				finalizers := object.GetFinalizers()
				Expect(finalizers).NotTo(ContainElement(PVCFinalizersForPod))
			}

			// cache should not be deleted
			err = fakeClient.Get(ctx, cachePVCKey, object)
			Expect(err).Should(BeNil())
			Expect(object.Name).Should(BeEquivalentTo("build-cache"))

		})

	})

})
