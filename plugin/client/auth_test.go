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

	"github.com/go-resty/resty/v2"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
)

func TestAuth_WithContext(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Type:   v1alpha1.AuthTypeBasic,
		Secret: map[string][]byte{"123": []byte("456")},
	}

	authCtx := auth.WithContext(context.TODO())

	g.Expect(authCtx).ToNot(BeNil())
}

func TestExtractAuth(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Type:   v1alpha1.AuthTypeBasic,
		Secret: map[string][]byte{"123": []byte("456")},
	}

	authCtx := auth.WithContext(context.TODO())
	newAuth := ExtractAuth(authCtx)

	g.Expect(newAuth).ToNot(BeNil())
	g.Expect(newAuth).To(Equal(auth))
}

func TestBasicAuth(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Type:   v1alpha1.AuthTypeBasic,
		Secret: map[string][]byte{"username": []byte("123"), "password": []byte("456")},
	}

	request := resty.New().R()
	err := auth.ToRequest(request)
	g.Expect(err).To(BeNil())
	g.Expect(request.UserInfo.Username).To(Equal("123"))
	g.Expect(request.UserInfo.Password).To(Equal("456"))

}

func TestOauth2(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Type:   v1alpha1.AuthTypeOauth2,
		Secret: map[string][]byte{"token": []byte("123")},
	}

	request := resty.New().R()
	err := auth.ToRequest(request)
	g.Expect(err).To(BeNil())
	g.Expect(request.Header.Get("Authorization")).To(Equal("Bearer 123"))

}

func TestPersonalToken(t *testing.T) {
	g := NewGomegaWithT(t)

	auth := &Auth{
		Type:   v1alpha1.AuthTypePersonalToken,
		Secret: map[string][]byte{"token": []byte("123")},
	}

	request := resty.New().R()
	err := auth.ToRequest(request)
	g.Expect(err).To(BeNil())
	g.Expect(request.Header.Get("Authorization")).To(Equal("Bearer 123"))
}
