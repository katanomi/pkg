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

package config

import (
	"fmt"
	"strconv"
	"time"
)

// FeatureValue definition of FeatureValue feature value
type FeatureValue string

// String returned as a string
func (f FeatureValue) String() string {
	return string(f)
}

// AsInt returned as an integer, or -1 if the conversion fails.
func (f FeatureValue) AsInt() (int, error) {
	v, err := strconv.Atoi(f.String())
	if err != nil {
		return -1, fmt.Errorf("failed parsing feature flags config %q: %v", f.String(), err)
	}
	return v, nil
}

// AsBool returns as a Duration, or 0 if the conversion fails.
func (f FeatureValue) AsDuration() (time.Duration, error) {
	v, err := time.ParseDuration(f.String())
	if err != nil {
		return 0, fmt.Errorf("failed parsing feature flags config %q: %v", f.String(), err)
	}
	return v, nil
}

// AsBool returns as a Bool, or false if the conversion fails.
func (f FeatureValue) AsBool() (bool, error) {
	v, err := strconv.ParseBool(f.String())
	if err != nil {
		return false, fmt.Errorf("failed parsing feature flags config %q: %v", f.String(), err)
	}
	return v, nil
}
