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

package url

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	trimSSHProtocol = regexp.MustCompile(`^.*?@`)
	hostPortFormat  = regexp.MustCompile(`[\d.]*:\d{1,5}`)
)

func UrlToProjectID(gitURL string) (string, error) {
	if strings.HasPrefix(gitURL, "https") || strings.HasPrefix(gitURL, "http") {
		return httpURLToProjectID(gitURL)
	}

	return sshURLToProjectID(gitURL)
}

// sshURLToProjectID converts ssh url to project id
// some common cases:
// ssh://git@192.168.130.62:31211/katanomi/catalog.git
// git@github.com:katanomi/catalog.git
// ssh://git@github.com/katanomi/catalog.git
// ssh://git@[2004::192:168:139:4]:32078/root/image.git
// git@[2004::192:168:139:4]:32078/root/image.git
// git@[2004::192:168:139:4]/root/image.git

func sshURLToProjectID(gitURL string) (string, error) {
	s := trimSSHProtocol.ReplaceAllString(gitURL, "")
	s = strings.TrimSuffix(s, ".git")
	if hostPortFormat.MatchString(s) {
		list := strings.SplitN(s, "/", 2)
		if len(list) == 2 {
			repoPath := strings.ReplaceAll(list[1], "/", "-")
			host := strings.ReplaceAll(list[0], "[", "")
			host = strings.ReplaceAll(host, "]", "")
			return host + "-" + strings.Trim(repoPath, "-"), nil
		}
	} else {
		s = strings.ReplaceAll(s, "/", "-")
		s = strings.ReplaceAll(s, ":", "-")
		s = strings.ReplaceAll(s, "[", "-")
		s = strings.ReplaceAll(s, "]", "-")
		return strings.Trim(s, "-"), nil
	}

	return gitURL, nil
}

// httpURLToProjectID converts http url to project id
func httpURLToProjectID(gitURL string) (string, error) {
	u, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}
	urlPath := strings.TrimSuffix(u.Path, ".git")
	urlPath = strings.ReplaceAll(urlPath, "/", "-")
	urlHost := u.Host
	urlHost = strings.ReplaceAll(urlHost, "[", "-")
	urlHost = strings.ReplaceAll(urlHost, "]", "-")

	return strings.Trim(urlHost, "-") + "-" + strings.Trim(urlPath, "-"), nil
}
