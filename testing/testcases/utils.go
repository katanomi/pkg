//go:build e2e
// +build e2e

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

package testcases

import (
	"context"
	"fmt"

	"github.com/katanomi/pkg/plugin/client"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var optionalGomega = gomega.NewGomega(func(message string, callerSkip ...int) {
	GinkgoWriter.Println("Assert failed:" + message)
})

func OptionalExpect(actual interface{}, extra ...interface{}) gomega.Assertion {
	return optionalGomega.Expect(actual, extra...)
}

// NameGetter get metadata name
type NameGetter interface {
	GetName() string
}

// ToPtrList convert a list to a list of pointer
func ToPtrList[T any](list []T) []*T {
	ptrList := make([]*T, 0, len(list))
	for index := range list {
		ptrList = append(ptrList, &list[index])
	}
	return ptrList
}

// FindByName find an item by name
func FindByName[T NameGetter](list []T, name string) T {
	index := FindIndexByName(list, name)
	if index == -1 {
		// FIXME：find a better way to return T(nil)
		var t T
		return t
	}
	return list[index]
}

// FindIndexByName find an item index by name
func FindIndexByName[T NameGetter](list []T, name string) int {
	for index, item := range list {
		if item.GetName() == name {
			return index
		}
	}
	return -1
}

func GetUsernamePasswordFromCtx(ctx context.Context) (username string, password string) {
	auth := client.ExtractAuth(ctx)
	if auth == nil {
		panic("no auth found")
	}

	username, password, err := auth.GetBasicInfo()
	if err != nil {
		panic(fmt.Sprintf("get basic auth info failed:%s", err))
	}
	return username, password
}
