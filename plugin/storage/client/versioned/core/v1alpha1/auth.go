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

	"github.com/katanomi/pkg/apis/storage/v1alpha1"
	client2 "github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/storage/client"
)

type AuthGetter interface {
	Auth() AuthInterface
}

type AuthInterface interface {
	Check(ctx context.Context, request v1alpha1.StorageAuthCheckRequest) (*v1alpha1.StorageAuthCheck, error)
}

type auth struct {
	client client.Interface
}

func (a *auth) Check(ctx context.Context, request v1alpha1.StorageAuthCheckRequest) (authCheck *v1alpha1.
	StorageAuthCheck,
	err error) {
	path := "auth/check"
	authCheck = &v1alpha1.StorageAuthCheck{}
	err = a.client.Post(ctx, path, client2.ResultOpts(authCheck), client2.BodyOpts(request))
	if err != nil {
		return nil, err
	}
	return
}

// newAuth returns an auth
func newAuth(c *CoreV1alpha1Client) *auth {
	return &auth{
		client: c.RESTClient(),
	}
}
