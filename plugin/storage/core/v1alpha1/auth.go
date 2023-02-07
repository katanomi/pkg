/*
Copyright 2023 The Katanomi Authors.

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

	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	apistoragev1alpha1 "github.com/katanomi/pkg/apis/storage/v1alpha1"
)

// AuthChecker checks auth according to params values
type AuthChecker interface {
	// CheckAuth used for auth checking of storage plugins
	CheckAuth(ctx context.Context, params []v1alpha1.Param) (*apistoragev1alpha1.StorageAuthCheck, error)
}
