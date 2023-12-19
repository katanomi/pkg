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

package base

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsNotImplementedError(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "normal error",
			err:  errors.New("normal error"),
			want: false,
		},
		{
			name: "not found error",
			err: &kerrors.StatusError{
				ErrStatus: metav1.Status{Code: http.StatusNotFound},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Expect(IsNotImplementedError(tt.err)).To(Equal(tt.want))
		})
	}
}
