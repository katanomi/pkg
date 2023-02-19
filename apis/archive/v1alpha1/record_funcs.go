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
	"github.com/katanomi/pkg/encoding"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ObjectToRecord convert k8s resource to a Record
func ObjectToRecord(obj client.Object) Record {
	data := unstructured.Unstructured{Object: encoding.ObjectToMap(obj)}
	// remove manage fields
	data.SetManagedFields(nil)
	data.SetDeletionTimestamp(nil)

	gvk := obj.GetObjectKind().GroupVersionKind()
	return Record{
		Spec: RecordSpec{
			UID:               string(obj.GetUID()),
			Group:             gvk.Group,
			Version:           gvk.Version,
			Kind:              gvk.Kind,
			Namespace:         obj.GetNamespace(),
			Name:              obj.GetName(),
			Data:              data.Object,
			CreationTimestamp: data.GetCreationTimestamp().Time.Unix(),
		},
	}
}

// TopConditionToMetadata convert top level condition to metadata
func TopConditionToMetadata(conds v1beta1.Conditions) map[string]string {
	for _, cond := range conds {
		if cond.Type == apis.ConditionSucceeded {
			metadata := map[string]string{}
			metadata["status"] = string(cond.Status)
			metadata["reason"] = cond.Reason
			return metadata
		}
	}
	return nil
}
