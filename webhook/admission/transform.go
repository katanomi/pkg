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
	admissionv1 "k8s.io/api/admission/v1"
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
		if triggeredBy.TriggeredTimestamp.IsZero() {
			creation := metaobj.GetCreationTimestamp()
			// if creation is not set, we set it to the current time.
			if creation.IsZero() {
				creation = metav1.Now()
			}
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

// WithCreatedBy adds a createdBy annotation to the object using the request information
// when an object already has the createdBy annotation it will only increment missing data
func WithCreatedBy() TransformFunc {
	return func(ctx context.Context, obj runtime.Object, req admission.Request) {
		if req.Operation != admissionv1.Create {
			return
		}
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
		createdBy := &mv1alpha1.CreatedBy{}
		createdBy, err = createdBy.FromAnnotation(annotations)
		if err != nil {
			log.Warnw("cannot unmarshal annotation value into createdBy struct", "err", err)
		}
		if createdBy == nil {
			createdBy = &mv1alpha1.CreatedBy{}
		}

		if createdBy.User == nil || createdBy.User.Name == "" {
			createdBy.User = SubjectFromRequest(req)
		}
		annotations, err = createdBy.SetIntoAnnotation(annotations)
		if err != nil {
			log.Warnw("cannot marshal createdBy struct to json ", "err", err, "struct", createdBy)
		} else {
			metaobj.SetAnnotations(annotations)
		}
	}
}
