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

package v1alpha1

import (
	"context"

	"github.com/katanomi/pkg/webhook/admission"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/logging"
)

//+kubebuilder:webhook:path=/mutate-test-katanomi-dev-v1alpha1-foobar,mutating=true,failurePolicy=fail,groups=test.katanomi.dev,resources=foobars,verbs=create;update,versions=v1alpha1,name=vfoobar.kb.io,sideEffects=None,admissionReviewVersions={v1,v1beta1}

var _ admission.Defaulter = &FooBar{}

// Default add default values for IntegationClass during create and update
func (i *FooBar) Default(ctx context.Context) {
	log := logging.FromContext(ctx)
	log.Debugw("running Default method for FooBar")
}

//+kubebuilder:webhook:path=/validate-test-katanomi-dev-v1alpha1-foobar,mutating=false,failurePolicy=fail,sideEffects=None,groups=test.katanomi.dev,resources=foobars,verbs=create;update,versions=v1alpha1,name=vfoobar.kb.io,admissionReviewVersions={v1,v1beta1}

var _ admission.Validator = &FooBar{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (i *FooBar) ValidateCreate(ctx context.Context) error {
	log := logging.FromContext(ctx)
	log.Debugw("running ValidateCreate method for FooBar")
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (i *FooBar) ValidateUpdate(ctx context.Context, old runtime.Object) error {
	log := logging.FromContext(ctx)
	log.Debugw("running ValidateUpdate method for FooBar")
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (i *FooBar) ValidateDelete(ctx context.Context) error {
	log := logging.FromContext(ctx)
	log.Debugw("running ValidateDelete method for FooBar")
	return nil
}
