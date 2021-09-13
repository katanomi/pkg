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
	"context"

	mv1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// TransformFuncused to make common defaulting logic amongst multiple resource
// using a context, an object and a request
type TransformFunc func(context.Context, runtime.Object, admission.Request)

// WithTriggeredBy adds a triggeredBy annotation to the object using the request information
// when an object already has the triggeredBy annotation it will only increment missing data
func WithTriggeredBy() TransformFunc {
	return func(ctx context.Context, obj runtime.Object, req admission.Request) {
		metaobj, ok := obj.(metav1.Object)
		if !ok {
			return
		}
		log := logging.FromContext(ctx)
		annotations := metaobj.GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
		}

		var err error
		triggeredBy := &mv1alpha1.TriggeredBy{}
		triggeredBy, err = triggeredBy.FromAnnotation(annotations)
		if err != nil {
			log.Warnw("cannot unmarshal annotation value into triggeredBy struct", "err", err)
		}
		if triggeredBy == nil {
			triggeredBy = &mv1alpha1.TriggeredBy{}
		}

		if triggeredBy.User == nil || triggeredBy.User.Name == "" {
			triggeredBy.User = SubjectFromRequest(req)
		}
		if triggeredBy.TriggeredTimestamp == nil {
			creation := metaobj.GetCreationTimestamp()
			triggeredBy.TriggeredTimestamp = &creation
		}
		annotations, err = triggeredBy.SetIntoAnnotation(annotations)
		if err != nil {
			log.Warnw("cannot marshal triggeredBy struct to json ", "err", err, "struct", triggeredBy)
		} else {
			metaobj.SetAnnotations(annotations)
		}
	}
}
