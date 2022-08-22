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
	// "fmt"
	"strings"
)

// GetKeyValues returns the set of keys and values for one flag in a set of given arguments
// args should be all the arguments to be interpreted. i.e --arg1 key1=value1 key2=value2 --arg2 key3=value3
// flag should be the flag name to select keys and values from, i.e arg1 or arg2 in the example above
func GetKeyValues(ctx context.Context, args []string, flag string) (keyValues map[string]string, ok bool) {
	keyValues = map[string]string{}
	var values []string
	if values, ok = GetArrayValues(ctx, args, flag); ok {
		for _, arg := range values {
			split := strings.SplitN(arg, "=", 2)
			if len(split) == 2 {
				keyValues[split[0]] = split[1]
			}
		}
	}
	return
}
