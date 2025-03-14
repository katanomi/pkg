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

package admission

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	mockadmission "github.com/katanomi/pkg/testing/mock/github.com/katanomi/pkg/webhook/admission"
	mockclient "github.com/katanomi/pkg/testing/mock/sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Context("Test.Approving.Handle", func() {
	DescribeTable("approval webhook",
		func(create bool, advanced bool, validatePassed bool, modifiedOthers bool, expected bool) {
			mockCtrl := gomock.NewController(GinkgoT())
			defer mockCtrl.Finish()

			mockClient := mockclient.NewMockClient(mockCtrl)
			mockClient.EXPECT().Scheme().Return(scheme.Scheme).AnyTimes()
			ctx := logging.WithLogger(context.Background(), logger)
			ctx = kclient.WithClient(ctx, mockClient)

			// Create a new webhook
			approving := generateApproving(ctx, mockCtrl, mockClient, advanced, validatePassed, modifiedOthers)

			req := baseRequest
			if create {
				req.Operation = admissionv1.Create
			} else {
				req.Operation = admissionv1.Update
			}

			// Test the handle method
			response := approving.Handle(ctx, req)
			Expect(response.Allowed).To(Equal(expected), "expected allowed to be %v", expected)
		},

		Entry("create-general-notpass-unmodifed-false", true, false, false, false, false),
		Entry("create-general-passed-modifed-true", true, false, true, true, true),
		Entry("create-general-passed-unmodifed-true", true, false, true, false, true),

		Entry("create-advanced-notpass-unmodifed-false", true, true, false, false, false),
		Entry("create-advanced-passed-modifed-true", true, true, true, true, true),
		Entry("create-advanced-passed-unmodifed-true", true, true, true, false, true),

		Entry("update-general-notpass-unmodifed-false", false, false, false, false, false),
		Entry("update-general-passed-modifed-false", false, false, true, true, false),
		Entry("update-general-passed-unmodifed-true", false, false, true, false, true),

		Entry("update-advanced-notpass-unmodifed-false", false, true, false, false, false),
		Entry("update-advanced-passed-modifed-true", false, true, true, true, true),
		Entry("update-advanced-passed-unmodifed-true", false, true, true, false, true),
	)
})

func validateApprovalPassed(ctx context.Context, reqUser authenticationv1.UserInfo, allowRepresentOthers, isCreateOperation bool,
	approvalSpecList []*metav1alpha1.ApprovalSpec, checkList []PairOfOldNewCheck, triggeredBy *metav1alpha1.TriggeredBy) (err error) {
	return nil
}

func validateApprovalRejected(ctx context.Context, reqUser authenticationv1.UserInfo, allowRepresentOthers, isCreateOperation bool,
	approvalSpecList []*metav1alpha1.ApprovalSpec, checkList []PairOfOldNewCheck, triggeredBy *metav1alpha1.TriggeredBy) (err error) {
	return fmt.Errorf("rejected")
}

func getResourceAttributes(verb string) authv1.ResourceAttributes {
	return authv1.ResourceAttributes{}
}

var (
	baseRequest = admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Operation: admissionv1.Update,
			Object: runtime.RawExtension{
				Raw: []byte(`{
	    "apiVersion": "v1",
	    "kind": "Pod",
	    "metadata": {
	        "name": "foo",
	        "namespace": "default"
	    },
	    "spec": { }
	}`),
			},
			OldObject: runtime.RawExtension{
				Raw: []byte(`{
	    "apiVersion": "v1",
	    "kind": "Pod",
	    "metadata": {
	        "name": "foo",
	        "namespace": "default"
	    },
	    "spec": { }
	}`),
			},
		},
	}
)

func generateApproving(ctx context.Context, mockCtrl *gomock.Controller, mockClient *mockclient.MockClient,
	advanced bool, validatePassed bool, modifiedOthers bool) *approvingHandler {
	// starts mock controller
	mockApproval := mockadmission.NewMockApprovalWithTriggeredByGetter(mockCtrl)
	mockApproval.EXPECT().DeepCopyObject().Return(&corev1.Pod{}).AnyTimes()
	mockApproval.EXPECT().GetApprovalSpecs(gomock.Any()).Return(nil).AnyTimes()
	mockApproval.EXPECT().GetChecks(gomock.Any()).Return(nil).AnyTimes()
	mockApproval.EXPECT().ModifiedOthers(gomock.Any(), gomock.Any()).Return(modifiedOthers).AnyTimes()
	mockApproval.EXPECT().GetTriggeredBy(gomock.Any()).Return(nil).AnyTimes()

	webhook := ApprovingWebhookFor(ctx, mockApproval, getResourceAttributes)
	approving, _ := webhook.Handler.(*approvingHandler)
	Expect(approving).NotTo(BeNil())

	approving.client = mockClient

	mockClient.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, obj client.Object, _ ...client.CreateOption) error {
			// update the .status.allowed field
			stype := reflect.ValueOf(obj).Elem()
			status := stype.FieldByName("Status")
			if status.IsValid() {
				allowed := status.FieldByName("Allowed")
				if allowed.IsValid() {
					allowed.SetBool(advanced)
				}
			}
			return nil
		}).AnyTimes()

	if validatePassed {
		approving.validateApproval = validateApprovalPassed
	} else {
		approving.validateApproval = validateApprovalRejected
	}

	return approving
}
