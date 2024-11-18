/*
Copyright 2022 The AlaudaDevops Authors.

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

package logger

import (
	"context"
	"fmt"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// ResultErrors print the result of validate
func ResultErrors(ctx context.Context, errs utilerrors.Aggregate, successMsg, errorMsg string) (err error) {
	log := NewLoggerFromContext(ctx)
	if errs == nil {
		log.Infof("==> ✅  %s", successMsg)
		return nil
	}

	Errors(ctx, errs)
	return fmt.Errorf("%d %s", len(errs.Errors()), errorMsg)
}

// Errors print the result of validate
func Errors(ctx context.Context, errs utilerrors.Aggregate) {
	if errs == nil {
		return
	}

	log := NewLoggerFromContext(ctx)
	for _, err := range errs.Errors() {
		log.Errorf("==> 🛑  %s", err)
	}
}
