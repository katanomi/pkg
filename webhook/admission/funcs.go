/*
Copyright 2021 The Katanomi Authors.

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
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// below both methods comes from sigs.k8s.io/controller-runtime/pkg/builder/webhook.go

func generateMutatePath(gvk schema.GroupVersionKind) string {
	return "/mutate-" + strings.Replace(gvk.Group, ".", "-", -1) + "-" +
		gvk.Version + "-" + strings.ToLower(gvk.Kind)
}

func generateValidatePath(gvk schema.GroupVersionKind) string {
	return "/validate-" + strings.Replace(gvk.Group, ".", "-", -1) + "-" +
		gvk.Version + "-" + strings.ToLower(gvk.Kind)
}

// returns a user based on the request information
func SubjectFromRequest(req admission.Request) *rbacv1.Subject {
	sub := &rbacv1.Subject{}
	if strings.HasPrefix(req.UserInfo.Username, "system:serviceaccount:") {
		sub.Kind = rbacv1.ServiceAccountKind
	} else {
		// all non-service accounts are treated as users
		sub.Kind = rbacv1.UserKind
	}
	sub.Name = req.UserInfo.Username
	return sub
}
