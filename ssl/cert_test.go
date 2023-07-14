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

package ssl

import (
	"encoding/asn1"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCanonicalString(t *testing.T) {
	cases := map[string][]struct {
		value    string
		expected string
	}{
		"trims leading and trailing whitespace": {
			{value: " some-val ", expected: "some-val"},
			{value: "\f\tsome-val\n\n\v", expected: "some-val"},
		},
		"replaces any remaining whitespace with a single ' '": {
			{value: "SOME  VAL", expected: "some val"},
			{value: "SOME\nVAL", expected: "some val"},
			{value: "SOME \f\t\n\r\vVAL", expected: "some val"},
		},
		`defines whitespaces as '\f' '\t' '\n' '\r'  and '\v' runes`: {
			{value: "\u0085some-val\u0085", expected: "\u0085some-val\u0085"},
			{value: "\u00A0some-val\u00A0", expected: "\u00A0some-val\u00A0"},
		},
		"converts to lowercase": {
			{value: "SOME VAL", expected: "some val"},
		},
	}

	for desc, values := range cases {
		t.Run(desc, func(t *testing.T) {
			for _, item := range values {
				g := NewGomegaWithT(t)
				g.Expect(CanonicalString(item.value)).To(Equal(item.expected))
			}
		})
	}
}

func TestCanonicalName(t *testing.T) {
	cases := []struct {
		path string
		cn   string
	}{
		{
			// generated with command
			// openssl genrsa -out ca.key 2048
			// openssl req -x509 -new -nodes -key ca.key -subj "/CN=<host>" -days 10000 -out <path>.crt
			path: "./testdata/abc.test.crt",
			// generated with command
			// openssl x509 -hash -fingerprint -noout -in <path>.crt
			cn: "abc.test",
		},
		{
			path: "./testdata/def.test.crt",
			cn:   "def.test",
		},
		{
			path: "./testdata/abc.dev.crt",
			cn:   "abc.dev",
		},
		{
			path: "./testdata/def.dev.crt",
			cn:   "def.dev",
		},
	}

	for i, item := range cases {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			g := NewGomegaWithT(t)

			raw, err := os.ReadFile(item.path)
			g.Expect(err).NotTo(HaveOccurred())

			cert, err := ParseCert(raw)
			g.Expect(err).NotTo(HaveOccurred())

			canonicalName, err := CanonicalName(cert.RawSubject)
			g.Expect(err).NotTo(HaveOccurred())

			set := canonicalSET{}
			_, err = asn1.Unmarshal(canonicalName, &set)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(set[0].Value).To(Equal(item.cn))
		})
	}
}

func TestCertHash(t *testing.T) {
	cases := []struct {
		path string
		hash string
	}{
		{
			// generated with command
			// openssl genrsa -out ca.key 2048
			// openssl req -x509 -new -nodes -key ca.key -subj "/CN=<host>" -days 10000 -out <path>.crt
			path: "./testdata/abc.test.crt",
			// generated with command
			// openssl x509 -hash -fingerprint -noout -in <path>.crt
			hash: "3e4f079c.0",
		},
		{
			path: "./testdata/def.test.crt",
			hash: "04002ec7.0",
		},
		{
			path: "./testdata/abc.dev.crt",
			hash: "53f61dcd.0",
		},
		{
			path: "./testdata/def.dev.crt",
			hash: "851a8cf9.0",
		},
	}

	for i, item := range cases {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			g := NewGomegaWithT(t)
			hash, err := CertHash(item.path)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(hash).To(Equal(item.hash))
		})
	}
}
