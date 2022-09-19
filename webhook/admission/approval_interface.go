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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

//go:generate mockgen -package=admission -destination=../../testing/mock/github.com/katanomi/pkg/webhook/admission/approval.go github.com/katanomi/pkg/webhook/admission Approval

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

//go:generate mockgen -package=admission -destination=../../testing/mock/github.com/katanomi/pkg/webhook/admission/triggeredbygetter.go github.com/katanomi/pkg/webhook/admission TriggeredByGetter

// TriggeredByGetter get the triggerd by from the runtime object
// This interface should be implemented when `requiresDifferentApprover` is enabled.
type TriggeredByGetter interface {
	GetTriggeredBy(runtime.Object) *metav1alpha1.TriggeredBy
}

//go:generate mockgen -package=admission -destination=../../testing/mock/github.com/katanomi/pkg/webhook/admission/approval_with_triggeredbygetter.go github.com/katanomi/pkg/webhook/admission ApprovalWithTriggeredByGetter

// ApprovalWithTriggeredByGetter defines functions for approving resources and enables `requiresDifferentApprover`
type ApprovalWithTriggeredByGetter interface {
	Approval
	TriggeredByGetter
}
