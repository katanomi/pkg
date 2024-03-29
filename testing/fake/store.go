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

package fake

import (
	"context"

	"github.com/katanomi/pkg/testing/fake/opa"
)

type Store interface {
	Setup(ctx context.Context) error
	Create(ctx context.Context, p *opa.Policy) error
	Get(ctx context.Context, id string) (*opa.Policy, error)
	Delete(ctx context.Context, id string) error
}
