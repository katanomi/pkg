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
	"fmt"

	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// ApprovalInterface is an interface for getting approval information from object.
type ApprovalInterface interface {
	ChecksGetter

	// GetApprovalSpecs returns the list of ApprovalSpecs for the given object.
	GetApprovalSpecs(runtime.Object) []*metav1alpha1.ApprovalSpec

	// GetDefault returns the default value of the resource.
	GetDefault() runtime.Object

	// ModifiedOthers returns true if the object has also modified other content.
	ModifiedOthers(runtime.Object, runtime.Object) bool
}

// ApprovalWebhook is used to verify that the approval is legitimate.
// +k8s:deepcopy-gen=false
type ApprovalWebhook struct {
	// ApprovalInterface used to get approval information from object.
	ApprovalInterface

	decoder *admission.Decoder
	Log     *zap.SugaredLogger
	Client  client.Client

	// ResourceAttributes used to determine if advanced permissions are available
	ResourceAttributes *authv1.ResourceAttributes
}

// Handle handles admission requests.
//  1. Review by local access determine if advanced permissions are available
//     Advanced permissions can approve on behalf of others and modify other content
//  2. Determine whether the approval operation is legal
//     Reject if not legal
//  3. Ordinary users are not allowed to modify other content
//     If modified, reject
func (a *ApprovalWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	new := a.GetDefault()
	old := a.GetDefault()

	a.decoder.DecodeRaw(req.Object, new)
	a.decoder.DecodeRaw(req.OldObject, old)

	obj := new.(client.Object)
	objectKey := fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
	log := a.Log.With("key", objectKey)
	ctx = logging.WithLogger(ctx, log)

	log.Debugw("approval webhook", "name", obj.GetName(), "user", req.AdmissionRequest.UserInfo,
		"old", req.AdmissionRequest.OldObject, "new", req.AdmissionRequest.Object)

	var extra map[string]authv1.ExtraValue
	if req.UserInfo.Extra != nil {
		extra = map[string]authv1.ExtraValue{}
		for k, v := range req.UserInfo.Extra {
			extra[k] = []string(v)
		}
	}

	resource := a.ResourceAttributes.DeepCopy()
	resource.Name = obj.GetName()
	resource.Namespace = obj.GetNamespace()
	localAccessReview := &authv1.LocalSubjectAccessReview{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
		},
		Spec: authv1.SubjectAccessReviewSpec{
			User:   req.UserInfo.Username,
			Groups: req.UserInfo.Groups,
			UID:    req.UserInfo.UID,
			Extra:  extra,

			ResourceAttributes: resource,
		},
	}

	err := a.Client.Create(ctx, localAccessReview)
	log.Debugw("LOCAL ACCESS REVIEW FOR UPDATE", "err", err, "status", localAccessReview.Status)
	advancedPermissions := localAccessReview.Status.Allowed

	isCreateOperation := (req.Operation == admissionv1.Create)
	approvalSpecList := a.GetApprovalSpecs(new)

	oldChecks := a.GetChecks(old)
	newChecks := a.GetChecks(new)
	oldNewChecksPairs := make([]PairOfOldNewCheck, len(oldChecks))
	for i := 0; i < len(oldChecks); i++ {
		oldNewChecksPairs[i] = PairOfOldNewCheck{oldChecks[i], newChecks[i]}
	}

	// Determine whether the approval act is legal
	if err := ValidateApproval(ctx, req.UserInfo, advancedPermissions, isCreateOperation, approvalSpecList, oldNewChecksPairs); err != nil {
		return admission.Denied(err.Error())
	}

	if advancedPermissions {
		return admission.Allowed("has advanced permissions")
	}

	// ordinary permissions can not patch other fields, except approval check
	if !isCreateOperation && a.ModifiedOthers(old, new) {
		log.Debugw("changed other fields", "old", old, "new", new)
		return admission.Denied("can not change other fields, except approval user input")
	}

	return admission.Allowed("allow approval")
}

func (a *ApprovalWebhook) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
