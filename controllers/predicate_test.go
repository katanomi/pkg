/*
Copyright 2024 The Katanomi Authors.

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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestSecretDataChangedPredicate(t *testing.T) {
	var data = []struct {
		desc string
		old  map[string][]byte
		new  map[string][]byte

		expected bool
	}{
		{
			desc: "old is nil",
			old:  nil,
			new: map[string][]byte{
				"a": []byte("1"),
			},
			expected: true,
		},
		{
			desc: "old is not nil and no changes",
			old: map[string][]byte{
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
			},
			expected: false,
		},
		{
			desc: "old is not nil and no changes 2",
			old: map[string][]byte{
				"b": []byte("0"),
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
				"b": []byte("0"),
			},
			expected: false,
		},
		{
			desc: "old is not nil and changes",
			old: map[string][]byte{
				"b": []byte("0"),
				"a": []byte("1"),
			},
			new: map[string][]byte{
				"a": []byte("1"),
				"b": []byte("1"),
			},
			expected: true,
		},
	}

	for _, item := range data {
		t.Run(item.desc, func(t *testing.T) {
			g := NewGomegaWithT(t)
			e := event.UpdateEvent{
				ObjectOld: &corev1.Secret{Data: item.old},
				ObjectNew: &corev1.Secret{Data: item.new},
			}
			actual := SecretDataChangedPredicate{}.Update(e)

			g.Expect(actual).Should(BeEquivalentTo(item.expected))
		})

	}
}
