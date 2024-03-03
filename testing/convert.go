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

package testing

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DefaultConvertRuntimeToClientobjectFunc convert runtime.Object to client.Object
// be similar to convertFromUnstructuredIfNecessary
func DefaultConvertRuntimeToClientobjectFunc(runtimeObj runtime.Object) (obj client.Object, err error) {
	obj, ok := runtimeObj.(client.Object)
	if !ok {
		err = fmt.Errorf("Unsupported gvk: %s", runtimeObj.GetObjectKind().GroupVersionKind())
	}
	return
}

// ConvertTypeMetaToGroupVersionResource converts type meta to group version resource
func ConvertTypeMetaToGroupVersionResource(typeMeta metav1.TypeMeta) schema.GroupVersionResource {
	gv, _ := schema.ParseGroupVersion(typeMeta.APIVersion)
	gvk := gv.WithKind(typeMeta.Kind)
	plural, _ := meta.UnsafeGuessKindToResource(gvk)
	return plural
}

// SliceToRuntimeOjbect convert slice to runtime.Object
func SliceToRuntimeOjbect[T any](s []T) []runtime.Object {
	r := make([]runtime.Object, 0, len(s))
	for _, v := range s {
		if o, ok := any(v).(runtime.Object); ok {
			r = append(r, o)
		}
	}
	if len(r) == 0 {
		return nil
	}
	return r
}

// SliceToInterfaceSlice convert a slice to a slice of interface
func SliceToInterfaceSlice[T any](s []T) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
