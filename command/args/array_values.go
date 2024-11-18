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

package args

import (
	"context"
	"fmt"
	"strings"
)

// GetArrayValues returns a list of values for a one flag in a set of given arguments
// args should be all the arguments to be interpreted. i.e ["--arg1","key1=value1","key2=value2", "--arg2", "key3=value3"]
// flag should be the flag name to select keys and values from, i.e arg1 or arg2 in the example above.
//
// Note: When the value is empty, it will be removed.
//
//	 i.e: flag: arg1
//		  args: ["--arg1", "", "key1=value1", "--arg2", "key3=value3"]
//	      return: ["key1=value1"]
func GetArrayValues(ctx context.Context, args []string, flag string) (values []string, ok bool) {
	values = make([]string, 0, len(args))
	offset := -1
	for i, arg := range args {
		if arg == fmt.Sprintf("--%s", flag) {
			// this is the key
			// save offset
			offset = i + 1
			ok = true
		} else if ok && strings.HasPrefix(arg, "--") {
			// found another flag and
			// already passed all keys, quit the loop
			values = args[offset:i]
			break
		} else if ok && i+1 == len(args) {
			// this is the last flag so needs to get it all
			values = args[offset : i+1]
		} else {
			// no-op
		}
	}

	// remove empty value
	i := 0
	for _, v := range values {
		if v != "" {
			values[i] = v
			i++
		}
	}
	values = values[0:i]

	return
}
