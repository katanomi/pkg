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
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestConfig_MustStringVal(t *testing.T) {
	g := NewGomegaWithT(t)

	tests := map[string]struct {
		m    DataMap
		key  string
		dft  string
		want string
	}{
		"config is nil, return default value": {
			m:    nil,
			key:  "key",
			dft:  "dft",
			want: "dft",
		},
		"config is empty, return default value": {
			m:    map[string]string{},
			key:  "key",
			dft:  "dft",
			want: "dft",
		},
		"config does not contains expected key, return default value": {
			m:    map[string]string{"not-exist-key": "not-exist-value"},
			key:  "key",
			dft:  "dft",
			want: "dft",
		},
		"config contains the expected key": {
			m:    map[string]string{"exist-key": "exist-value"},
			key:  "exist-key",
			dft:  "dft",
			want: "exist-value",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.m.MustStringVal(tt.key, tt.dft)).To(Equal(tt.want))
		})
	}
}

func TestConfig_MustIntVal(t *testing.T) {
	g := NewGomegaWithT(t)

	tests := map[string]struct {
		m    DataMap
		key  string
		dft  int
		want int
	}{
		"config is nil, return default value": {
			m:    nil,
			key:  "key",
			dft:  123,
			want: 123,
		},
		"config is empty, return default value": {
			m:    map[string]string{},
			key:  "key",
			dft:  123,
			want: 123,
		},
		"config does not contains expected key, return default value": {
			m:    map[string]string{"not-exist-key": "not-exist-value"},
			key:  "key",
			dft:  123,
			want: 123,
		},
		"config contains the expected key but the value is invalid": {
			m:    map[string]string{"exist-key": "123x"},
			key:  "exist-key",
			dft:  123,
			want: 123,
		},
		"config contains the expected key": {
			m:    map[string]string{"exist-key": "1234"},
			key:  "exist-key",
			dft:  123,
			want: 1234,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.m.MustIntVal(tt.key, tt.dft)).To(Equal(tt.want))
		})
	}
}

func TestConfig_MustBoolVal(t *testing.T) {
	g := NewGomegaWithT(t)

	tests := map[string]struct {
		m    DataMap
		key  string
		dft  bool
		want bool
	}{
		"config is nil, return default value": {
			m:    nil,
			key:  "key",
			dft:  true,
			want: true,
		},
		"config is empty, return default value": {
			m:    map[string]string{},
			key:  "key",
			dft:  true,
			want: true,
		},
		"config does not contains expected key, return default value": {
			m:    map[string]string{"not-exist-key": "not-exist-value"},
			key:  "key",
			dft:  true,
			want: true,
		},
		"config contains the expected key but the value is invalid": {
			m:    map[string]string{"exist-key": "falsex"},
			key:  "exist-key",
			dft:  true,
			want: true,
		},
		"config contains the expected key": {
			m:    map[string]string{"exist-key": "false"},
			key:  "exist-key",
			dft:  true,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.m.MustBoolVal(tt.key, tt.dft)).To(Equal(tt.want))
		})
	}
}

func TestConfig_MustTimeDurationVal(t *testing.T) {
	g := NewGomegaWithT(t)

	tests := map[string]struct {
		m    DataMap
		key  string
		dft  time.Duration
		want time.Duration
	}{
		"config is nil, return default value": {
			m:    nil,
			key:  "key",
			dft:  time.Second,
			want: time.Second,
		},
		"config is empty, return default value": {
			m:    map[string]string{},
			key:  "key",
			dft:  time.Second,
			want: time.Second,
		},
		"config does not contains expected key, return default value": {
			m:    map[string]string{"not-exist-key": "not-exist-value"},
			key:  "key",
			dft:  time.Second,
			want: time.Second,
		},
		"config contains the expected key but the value is invalid": {
			m:    map[string]string{"exist-key": "1ss"},
			key:  "exist-key",
			dft:  time.Second,
			want: time.Second,
		},
		"config contains the expected key": {
			m:    map[string]string{"exist-key": "1m"},
			key:  "exist-key",
			dft:  time.Second,
			want: time.Minute,
		},
		"config contains the expected key, and has day unit": {
			m:    map[string]string{"exist-key": "1d"},
			key:  "exist-key",
			dft:  time.Hour * 24,
			want: time.Hour * 24,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g.Expect(tt.m.MustTimeDurationVal(tt.key, tt.dft)).To(Equal(tt.want))
		})
	}
}
