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

package url

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"
)

// FullIPv6URL if the url contains an IPv6 address, convert compressed IPv6 to full format.
func FullIPv6URL(rawURL string) (string, error) {
	originURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if IsIPv6(originURL.Hostname()) {
		expandIPv6, _ := netip.ParseAddr(originURL.Hostname())
		if originURL.Port() != "" {
			originURL.Host = fmt.Sprintf("[%s]:%s", expandIPv6.StringExpanded(), originURL.Port())
		} else {
			originURL.Host = fmt.Sprintf("[%s]", expandIPv6.StringExpanded())
		}
	}
	return originURL.String(), nil
}

// IsIPv6 if the input is an IPv6 address, return true.
func IsIPv6(ipaddress string) bool {
	ip := net.ParseIP(ipaddress)
	return ip != nil && strings.Contains(ipaddress, ":")
}
