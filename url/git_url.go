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
	"net/netip"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"knative.dev/pkg/apis"
)

// MatchGitURLPrefix determine if the git URL is a subset of the target URL
// This matches:
// gitURL: https://github.com/katanomi/pkg.git
// target:
//   - https://github.com/katanomi/pkg.git
//   - https://github.com/katanomi/pkg
//   - https://github.com/katanomi/
//   - https://github.com
//
// This mismatch:
// gitURL: https://github.com/katanomi/pkg.git
// target:
//   - https://github.com/katanomi/pkg/
//
// This mismatch:
// gitURL: https://github.com/katanomi/pkg
// target:
//   - https://github.com/katanomi/pkg/
func MatchGitURLPrefix(gitURL, target *apis.URL) bool {
	if gitURL == nil || target == nil {
		return false
	}
	if gitURL.Host != target.Host {
		return false
	}
	sPath := gitURL.Path
	tPath := target.Path
	if strings.HasSuffix(sPath, ".git") {
		sPath = strings.TrimSuffix(sPath, ".git")
		tPath = strings.TrimSuffix(tPath, ".git")
	}
	if strings.HasSuffix(tPath, ".git") {
		// If the exact match, return true
		if sPath == strings.TrimSuffix(tPath, ".git") {
			return true
		}
	}
	// determine if the prefix matches
	return strings.HasPrefix(sPath, tPath)
}

// ExpandURLIPv6 if the url contains an IPv6 address, expand the IPv6 address.
func ExpandURLIPv6(rawURL string) (string, error) {
	originURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if govalidator.IsIPv6(originURL.Hostname()) {
		expandIPv6, _ := netip.ParseAddr(originURL.Hostname())
		if originURL.Port() != "" {
			originURL.Host = fmt.Sprintf("[%s]:%s", expandIPv6.StringExpanded(), originURL.Port())
		} else {
			originURL.Host = fmt.Sprintf("[%s]", expandIPv6.StringExpanded())
		}
	}
	return originURL.String(), nil
}
