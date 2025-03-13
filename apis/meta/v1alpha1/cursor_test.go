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
	"encoding/base64"
	"testing"

	. "github.com/onsi/gomega"
)

func TestTimeCursor_success(t *testing.T) {
	g := NewGomegaWithT(t)
	cursor := TimeCursor{
		Pager: Pager{
			ItemsPerPage: 30,
			Page:         2,
		},
		QueryStartAt: 1675650082,
	}
	expectStr := "eyJpdGVtc1BlclBhZ2UiOjMwLCJwYWdlIjoyLCJxdWVyeVN0YXJ0QXQiOjE2NzU2NTAwODJ9"
	g.Expect(cursor.Encode()).To(Equal(expectStr))

	timeCursor, err := ParseTimeCursor(expectStr)
	g.Expect(err).To(BeNil())
	g.Expect(timeCursor).To(Equal(&cursor))
}

func TestTimeCursor_fail(t *testing.T) {
	g := NewGomegaWithT(t)
	// invalid base64 string
	invalidBase64Str := "YWJjZA====="
	cursor, err := ParseTimeCursor(invalidBase64Str)
	g.Expect(cursor).To(BeNil())
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("illegal base64 data"))

	// invalid json content
	invalidContentStr := base64.StdEncoding.EncodeToString([]byte(`{"page":"1"}`))
	cursor, err = ParseTimeCursor(invalidContentStr)
	g.Expect(cursor).To(BeNil())
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("cannot unmarshal string into Go struct field TimeCursor.Pager.page of type int"))
}
