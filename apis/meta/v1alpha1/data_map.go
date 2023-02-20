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
	"strconv"
	"time"

	"github.com/k1LoW/duration"
)

// DataMap describe a map[string]string struct
type DataMap map[string]string

// MustStringVal return string value with specified key
func (c DataMap) MustStringVal(key string, dft string) string {
	s := c.StringVal(key)
	if s == nil {
		return dft
	}
	return *s
}

// StringVal return string value with specified key, if not found return nil
func (c DataMap) StringVal(key string) *string {
	if len(c) == 0 {
		return nil
	}
	s, ok := c[key]
	if !ok {
		return nil
	}
	return &s
}

// MustIntVal return int value with specified key
func (c DataMap) MustIntVal(key string, dft int) int {
	i, _ := c.IntVal(key)
	if i == nil {
		return dft
	}
	return *i
}

// IntVal return int value with specified key, if not found return nil
func (c DataMap) IntVal(key string) (*int, error) {
	v := c.StringVal(key)
	if v == nil {
		return nil, nil
	}
	i, err := strconv.Atoi(*v)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// MustBoolVal return bool value with specified key
func (c DataMap) MustBoolVal(key string, dft bool) bool {
	b, _ := c.BoolVal(key)
	if b == nil {
		return dft
	}
	return *b
}

// BoolVal return bool value with specified key, if not found return nil
func (c DataMap) BoolVal(key string) (*bool, error) {
	v := c.StringVal(key)
	if v == nil {
		return nil, nil
	}
	b, err := strconv.ParseBool(*v)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// MustTimeDurationVal return time.Duration value with specified key
func (c DataMap) MustTimeDurationVal(key string, dft time.Duration) time.Duration {
	d, _ := c.TimeDurationVal(key)
	if d == nil {
		return dft
	}
	return *d
}

// TimeDurationVal return time.Duration value with specified key, if not found return nil
func (c DataMap) TimeDurationVal(key string) (*time.Duration, error) {
	v := c.StringVal(key)
	if v == nil {
		return nil, nil
	}
	d, err := duration.Parse(*v)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
