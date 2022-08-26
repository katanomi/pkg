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

package approval

import (
	"context"

	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/namespace"
	"github.com/katanomi/pkg/user/matching"
	kadmission "github.com/katanomi/pkg/webhook/admission"
)

// PairOfOldNewCheck is a pair of old and new check
type PairOfOldNewCheck [2]*metav1alpha1.Check

// GetChecksFromObject gets the checks from the runtime object
type GetChecksFromObject func(ctx context.Context, obj runtime.Object) []*metav1alpha1.Check

// WithApprovalOperator adds an approval operator to the object using the request information
func WithApprovalOperator(getChecksFromObject GetChecksFromObject) kadmission.TransformFunc {
	return func(ctx context.Context, runtimeObj runtime.Object, req admission.Request) {
		if req.Operation != admissionv1.Create && req.Operation != admissionv1.Update {
			return
		}
		log := logging.FromContext(ctx)
		newChecks := getChecksFromObject(ctx, runtimeObj)
		var oldChecks []*metav1alpha1.Check
		if req.Operation == admissionv1.Create {
			// If it is a the create operation, the base is nil.
			oldChecks = make([]*metav1alpha1.Check, len(newChecks))
		} else {
			base := apis.GetBaseline(ctx)
			oldObject, _ := base.(runtime.Object)
			oldChecks = getChecksFromObject(ctx, oldObject)
		}
		if len(oldChecks) != len(newChecks) {
			log.Warnw("unable to add approval operator, length mismatch", "req", req, "old", oldChecks, "new", newChecks)
			return
		}

		// log.Debugw("add approval operator", "user", req.AdmissionRequest.UserInfo,
		// 	"old", req.AdmissionRequest.OldObject, "new", req.AdmissionRequest.Object)

		oldNewChecksPairs := make([]PairOfOldNewCheck, len(oldChecks))
		for i := 0; i < len(oldChecks); i++ {
			oldNewChecksPairs[i] = PairOfOldNewCheck{oldChecks[i], newChecks[i]}
		}

		ctx = namespace.WithNamespace(ctx, req.Namespace)
		// Assuming that administrator privileges already exist, add an approval operator information.
		// Finally, there have validation webhook to verify the legitimacy.
		addApprovalOperator(ctx, req.UserInfo, oldNewChecksPairs)
	}
}

// addApprovalOperator adds an real approval operator to the new check
func addApprovalOperator(ctx context.Context, reqUser authenticationv1.UserInfo, checksList []PairOfOldNewCheck) {
	log := logging.FromContext(ctx)
	for _, checks := range checksList {
		oldUsers := checks[0].GetApprovalUsers()
		newUsers := checks[1].GetApprovalUsers()

		// If the approver is someone else, add the real operator of the approval.
		for i, newUser := range newUsers {
			oldUser := oldUsers.GetBySubject(newUser.Subject)
			if (oldUser == nil || oldUser.Input == nil) && newUser.Input != nil {
				if !matching.IsRightUser(reqUser, newUser.Subject) {
					ns, _ := namespace.NamespaceFrom(ctx)
					subject := matching.ConvertUserInfoToSubject(reqUser, ns)
					newUsers[i].Operator = &subject
					log.Debugw("approval for others", "approver", newUser.Subject.Name, "operator", subject)
				} else if newUsers[i].Operator != nil {
					log.Warnw("clear approval operator", "approver", newUser.Subject.Name, "operator", newUser.Operator)
					newUsers[i].Operator = nil
				} else {
					// Normal approval without modification
				}
			}
		}
	}
}
