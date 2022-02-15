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
	"fmt"
	"reflect"
	"strings"

	"github.com/katanomi/pkg/sharedmain"
	"knative.dev/pkg/logging"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ValidatorWebhook interface {
	Validator
	sharedmain.WebhookSetup
	sharedmain.WebhookRegisterSetup
	WithValidateCreate(creates ...ValidateCreateFunc) ValidatorWebhook
	WithValidateUpdate(updates ...ValidateUpdateFunc) ValidatorWebhook
	WithValidateDelete(deletes ...ValidateDeleteFunc) ValidatorWebhook
}

type validatorWebhook struct {
	Validator
	creates []ValidateCreateFunc
	updates []ValidateUpdateFunc
	deletes []ValidateDeleteFunc
}

func (d *validatorWebhook) WithValidateCreate(creates ...ValidateCreateFunc) ValidatorWebhook {
	d.creates = creates
	return d
}

func (d *validatorWebhook) WithValidateUpdate(updates ...ValidateUpdateFunc) ValidatorWebhook {
	d.updates = updates
	return d
}

func (d *validatorWebhook) WithValidateDelete(deletes ...ValidateDeleteFunc) ValidatorWebhook {
	d.deletes = deletes
	return d
}

func NewValidatorWebhook(validator Validator) ValidatorWebhook {
	return &validatorWebhook{
		Validator: validator,
	}
}

func (d *validatorWebhook) GetLoggerName() string {
	if d.Validator != nil {
		typeName := strings.ToLower(reflect.TypeOf(d.Validator).Name())

		return fmt.Sprintf("%s-webhook-validation", typeName)
	}

	return "webhook-validation"
}

func (r *validatorWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r.Validator).
		Complete()
}

func (d *validatorWebhook) SetupRegisterWithManager(ctx context.Context, mgr ctrl.Manager) {
	log := logging.FromContext(ctx)

	if d.Validator == nil {
		log.Fatalw("webhook validator required")
		return
	}

	err := RegisterValidateWebhookFor(ctx, mgr, d.Validator, d.creates, d.updates, d.deletes)
	if err != nil {
		log.Fatalw("register webhook failed", "err", err)
	}
}
