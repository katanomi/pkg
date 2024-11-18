/*
Copyright 2022 The AlaudaDevops Authors.

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

package regex

import (
	"fmt"
	"regexp"
)

// regexes provide helper functions for regular expression array
type regexes struct {
	items   []string
	regexes []*regexp.Regexp
}

// Regexes will parse regular expression string array to regexes struct
func Regexes(items []string) *regexes {
	return &regexes{items: items, regexes: make([]*regexp.Regexp, 0, len(items))}
}

// Compile will compile all regular expressions
func (r *regexes) Compile() error {
	if len(r.items) == 0 {
		return nil
	}

	for _, item := range r.items {
		regex, err := regexp.Compile(item)
		if err != nil {
			return fmt.Errorf("regex '%s' compile error: %s", item, err.Error())
		}
		r.regexes = append(r.regexes, regex)
	}

	return nil
}

// MatchString will match string in any regular expression string array
func (r *regexes) MatchString(s string) (bool, error) {
	err := r.Compile()
	if err != nil {
		return false, err
	}
	return r.matchString(s), nil
}

func (r *regexes) matchString(s string) bool {
	if len(r.regexes) == 0 {
		return false
	}

	for _, regex := range r.regexes {
		if regex.MatchString(s) {
			return true
		}
	}

	return false
}

// MatchAnyString will match any string in any regular expression string array
func (r *regexes) MatchAnyString(sItems ...string) (matchedStrings []string, err error) {
	empty := []string{}

	err = r.Compile()
	if err != nil {
		return empty, err
	}

	matchedStrings = []string{}
	for _, s := range sItems {
		if r.matchString(s) {
			matchedStrings = append(matchedStrings, s)
		}
	}

	return matchedStrings, nil
}
