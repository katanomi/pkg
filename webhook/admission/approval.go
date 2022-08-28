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
	"net/http"

	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
)

//go:generate mockgen -package=admission -destination=../../testing/mock/github.com/katanomi/pkg/webhook/admission/approval.go  github.com/katanomi/pkg/webhook/admission Approval

// Approval defines functions for approving resources
type Approval interface {
	runtime.Object
	metav1.Object

	// ChecksGetter gets the checks from the runtime object
	ChecksGetter

	// GetApprovalSpecs returns the list of ApprovalSpecs for the given object.
	// Used to determine if advanced permissions are available
	GetApprovalSpecs(runtime.Object) []*metav1alpha1.ApprovalSpec

	// ModifiedOthers returns true if the object has also modified other content.
	ModifiedOthers(runtime.Object, runtime.Object) bool
}

// GetResourceAttributes returns the specified verb of resouce attributes.
type GetResourceAttributes func(string) authv1.ResourceAttributes

// ApprovingWebhookFor creates a new Webhook for Approving the provided type.
func ApprovingWebhookFor(ctx context.Context, approval Approval, getResourceAttributes GetResourceAttributes) *admission.Webhook {
	client := kclient.Client(ctx)
	if client == nil {
		panic("approval webhook need client")
	}
	return &admission.Webhook{
		Handler: &approvingHandler{
			approval:              approval,
			client:                client,
			getResourceAttributes: getResourceAttributes,
			SugaredLogger:         logging.FromContext(ctx),
			validateApproval:      ValidateApproval,
		},
	}
}

// approvingHandler is used to verify that the approval is legitimate.
// +k8s:deepcopy-gen=false
type approvingHandler struct {
	approval Approval
	decoder  *admission.Decoder

	client                client.Client
	getResourceAttributes GetResourceAttributes

	*zap.SugaredLogger

	validateApproval ValidateApprovalFunc
}

var _ admission.DecoderInjector = &approvingHandler{}

// InjectDecoder injects the decoder into a approvingHandler.
func (h *approvingHandler) InjectDecoder(d *admission.Decoder) error {
	h.decoder = d
	return nil
}

// Handle handles admission requests.
//  1. Review by local access determine if advanced permissions are available
//     Advanced permissions can approve on behalf of others and modify other content
//  2. Determine whether the approval operation is legal
//     Reject if not legal
//  3. General users are not allowed to modify other content
//     If modified, reject
func (h *approvingHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	if req.Operation != admissionv1.Create && req.Operation != admissionv1.Update {
		return admission.Allowed("only create and update operations are supported")
	}
	if h.approval == nil {
		panic("approval should never be nil")
	}

	// Get the object in the request
	old := h.approval.DeepCopyObject() //.(Approval)
	new := h.approval.DeepCopyObject() //.(Approval)

	err := h.decoder.DecodeRaw(req.Object, new)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if req.Operation == admissionv1.Update {
		err = h.decoder.DecodeRaw(req.OldObject, old)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	obj, _ := new.(client.Object)
	objectKey := fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
	log := h.SugaredLogger.With("key", objectKey)
	ctx = logging.WithLogger(ctx, log)

	log.Debugw("approval handle", "name", obj.GetName(), "user", req.AdmissionRequest.UserInfo,
		"old", req.AdmissionRequest.OldObject, "new", req.AdmissionRequest.Object)

	var extra map[string]authv1.ExtraValue
	if req.UserInfo.Extra != nil {
		extra = map[string]authv1.ExtraValue{}
		for k, v := range req.UserInfo.Extra {
			extra[k] = []string(v)
		}
	}

	resource := h.getResourceAttributes("update")
	resource.Name = obj.GetName()
	resource.Namespace = obj.GetNamespace()
	//
	localAccessReview := &authv1.LocalSubjectAccessReview{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: req.Namespace,
		},
		Spec: authv1.SubjectAccessReviewSpec{
			User:   req.UserInfo.Username,
			Groups: req.UserInfo.Groups,
			UID:    req.UserInfo.UID,
			Extra:  extra,

			ResourceAttributes: &resource,
		},
	}

	err = h.client.Create(ctx, localAccessReview, &client.CreateOptions{})
	if err != nil {
		log.Warnw("LOCAL ACCESS REVIEW FOR UPDATE", "resource", resource, "error", err)
		return admission.Denied(err.Error())
	}
	log.Debugw("LOCAL ACCESS REVIEW FOR UPDATE", "resource", resource, "status", localAccessReview.Status)
	advancedPermissions := localAccessReview.Status.Allowed

	isCreateOperation := (req.Operation == admissionv1.Create)
	approvalSpecList := h.approval.GetApprovalSpecs(new)

	oldChecks := h.approval.GetChecks(old)
	newChecks := h.approval.GetChecks(new)
	oldNewChecksPairs := make([]PairOfOldNewCheck, len(oldChecks))
	for i := 0; i < len(oldChecks); i++ {
		oldNewChecksPairs[i] = PairOfOldNewCheck{oldChecks[i], newChecks[i]}
	}

	// Determine whether the approval act is legal
	err = h.validateApproval(ctx, req.UserInfo, advancedPermissions, isCreateOperation, approvalSpecList, oldNewChecksPairs)
	if err != nil {
		return admission.Denied(err.Error())
	}

	if advancedPermissions {
		return admission.Allowed("has advanced permissions")
	}

	// General permissions can not patch other fields, except approval check
	if !isCreateOperation && h.approval.ModifiedOthers(old, new) {
		log.Debugw("changed other fields", "old", old, "new", new)
		return admission.Denied("can not change other fields, except approval user input")
	}

	return admission.Allowed("allow approval")
}
