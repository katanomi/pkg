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

package args

import (
	"context"
	"fmt"
	// "fmt"
	"strings"

	"github.com/katanomi/pkg/common"
	"k8s.io/apimachinery/pkg/util/errors"
)

// ValuesValidateOption provides extensive func to set validation rules of GetKeyValues
type ValuesValidateOption func(values []string) error

// ValuesValidationOptRequired is the builtin option for validating values are not empty
var ValuesValidationOptRequired ValuesValidateOption = func(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("empty values")
	}
	return nil
}

// ValuesValidationOptDuplicatedKeys is the builtin option for validating values are unique key=value format
var ValuesValidationOptDuplicatedKeys ValuesValidateOption = func(pairs []string) error {
	var duplicatedKeys, invalidPairs []string
	keyCount := map[string]int{}

	for _, pair := range pairs {
		split := strings.SplitN(pair, "=", 2)
		if len(split) != 2 && !common.Contains(invalidPairs, pair) {
			invalidPairs = append(invalidPairs, pair)
			continue
		}
		pairKey := split[0]
		keyCount[pairKey]++
		if keyCount[pairKey] == 2 {
			duplicatedKeys = append(duplicatedKeys, pairKey)
		}
	}

	if len(duplicatedKeys) == 0 && len(invalidPairs) == 0 {
		return nil
	}

	var duplicatedKeysErrMsg, invalidPairsErrMsg string
	if len(duplicatedKeys) > 0 {
		duplicatedKeysErrMsg = fmt.Sprintf("duplicated env keys: %s", duplicatedKeys)
	}

	if len(invalidPairs) > 0 {
		invalidPairsErrMsg = fmt.Sprintf("invalid env items: %s", invalidPairs)
	}
	return fmt.Errorf("invalid env-flags: %s", strings.Trim(strings.Join([]string{duplicatedKeysErrMsg,
		invalidPairsErrMsg},
		","), ","))
}

// GetKeyValues returns the set of keys and values for one flag in a set of given arguments
// args should be all the arguments to be interpreted. i.e --arg1 key1=value1 key2=value2 --arg2 key3=value3
// flag should be the flag name to select keys and values from, i.e arg1 or arg2 in the example above
func GetKeyValues(ctx context.Context, args []string, flag string,
	opts ...ValuesValidateOption) (keyValues map[string]string,
	err error) {
	keyValues = map[string]string{}
	values, _ := GetArrayValues(ctx, args, flag)
	for _, arg := range values {
		split := strings.SplitN(arg, "=", 2)
		if len(split) == 2 {
			keyValues[split[0]] = split[1]
		}
	}
	var errs []error
	for _, opt := range opts {
		if optErr := opt(values); optErr != nil {
			errs = append(errs, optErr)
		}
	}
	return keyValues, errors.NewAggregate(errs)
}
