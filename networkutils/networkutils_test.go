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

// Ref: https://github.com/openshift/library-go/blob/d679fe6f824818b04acb36524917c7362de6b81e/pkg/network/networkutils/networkutils_test.go

package networkutils

import (
	"strings"
	"testing"
)

func TestParseCIDRMask(t *testing.T) {
	tests := []struct {
		cidr       string
		fixedShort string
		fixedLong  string
	}{
		{
			cidr: "192.168.0.0/16",
		},
		{
			cidr: "192.168.1.0/24",
		},
		{
			cidr: "192.168.1.1/32",
		},
		{
			cidr:       "192.168.1.0/16",
			fixedShort: "192.168.0.0/16",
			fixedLong:  "192.168.1.0/32",
		},
		{
			cidr:       "192.168.1.1/24",
			fixedShort: "192.168.1.0/24",
			fixedLong:  "192.168.1.1/32",
		},
	}

	for _, test := range tests {
		_, err := ParseCIDRMask(test.cidr)
		if test.fixedShort == "" && test.fixedLong == "" {
			if err != nil {
				t.Fatalf("unexpected error parsing CIDR mask %q: %v", test.cidr, err)
			}
		} else {
			if err == nil {
				t.Fatalf("unexpected lack of error parsing CIDR mask %q", test.cidr)
			}
			if !strings.Contains(err.Error(), test.fixedShort) {
				t.Fatalf("error does not contain expected string %q: %v", test.fixedShort, err)
			}
			if !strings.Contains(err.Error(), test.fixedLong) {
				t.Fatalf("error does not contain expected string %q: %v", test.fixedLong, err)
			}
		}
	}
}

func TestIsPrivateAddress(t *testing.T) {
	for _, tc := range []struct {
		address string
		isLocal bool
	}{
		{"localhost", true},
		{"example.com", false},
		{"registry.localhost", false},

		{"9.255.255.255", false},
		{"10.0.0.1", true},
		{"10.1.255.255", true},
		{"10.255.255.255", true},
		{"11.0.0.1", false},

		{"127.0.0.1", true},

		{"172.15.255.253", false},
		{"172.16.0.1", true},
		{"172.30.0.1", true},
		{"172.31.255.255", true},
		{"172.32.0.1", false},

		{"192.167.122.1", false},
		{"192.168.0.1", true},
		{"192.168.122.1", true},
		{"192.168.255.255", true},
		{"192.169.1.1", false},

		{"::1", true},

		{"fe00::1", false},
		{"fd12:3456:789a:1::1", true},
		{"fe82:3456:789a:1::1", true},
		{"ff00::1", false},
	} {
		res := IsPrivateAddress(tc.address)
		if tc.isLocal && !res {
			t.Errorf("address %q considered not local", tc.address)
			continue
		}
		if !tc.isLocal && res {
			t.Errorf("address %q considered local", tc.address)
		}
	}
}
