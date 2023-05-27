/*
Copyright 2022 The Katanomi Authors.

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

// Package user package is to manage the context of user information
package user

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

type userInfoKey struct{}

// WithUserInfo returns a copy of parent in which the userinfo value is set
func WithUserInfo(parent context.Context, userinfo metav1alpha1.UserInfo) context.Context {
	return context.WithValue(parent, userInfoKey{}, userinfo)
}

// UserInfoFrom returns the value of the userinfo key on the ctx
func UserInfoFrom(ctx context.Context) (metav1alpha1.UserInfo, bool) {
	userinfo, ok := ctx.Value(userInfoKey{}).(metav1alpha1.UserInfo)
	return userinfo, ok
}

// UserInfoValue returns the value of the userinfo key on the ctx, or the empty string if none
func UserInfoValue(ctx context.Context) (result metav1alpha1.UserInfo) {
	userinfo, _ := UserInfoFrom(ctx)
	return userinfo
}

var entityInCxtKey = struct{}{}

// WithEntity will save entity into context
// up to now, we will use it to save operation target object in UserOwnedResouorceFilter
// and you could use EntityFromContext to avoid reading and unmarshall data from request again
func WithEntity(ctx context.Context, entity *unstructured.Unstructured) context.Context {
	ctx = context.WithValue(ctx, entityInCxtKey, entity)
	return ctx
}

// EntityFromContext will read entity from context
// up to now, we will use it to save operation target object in UserOwnedResouorceFilter
// and you could use EntityFromContext to avoid reading and unmarshall data from request again
func EntityFromContext(ctx context.Context) *unstructured.Unstructured {
	entity := ctx.Value(entityInCxtKey)
	return entity.(*unstructured.Unstructured)
}
