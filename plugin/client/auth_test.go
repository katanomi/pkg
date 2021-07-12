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

package client

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAuth_WithContext(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Method: "basic",
		Secret: "123",
	}

	authCtx := auth.WithContext(context.TODO())

	g.Expect(authCtx).ToNot(BeNil())
}

func TestExtractAuth(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Method: "basic",
		Secret: "123",
	}

	authCtx := auth.WithContext(context.TODO())
	newAuth := ExtractAuth(authCtx)

	g.Expect(newAuth).ToNot(BeNil())
	g.Expect(newAuth).To(Equal(auth))
}
