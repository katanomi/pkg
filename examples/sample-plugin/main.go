/*
Copyright 2021 The Katanomi Authors.

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

package main

import (
	"fmt"
	"os"

	"github.com/katanomi/pkg/examples/sample-plugin/harbor"

	"github.com/katanomi/pkg/plugin"
)

func main() {
	err := plugin.NewPlugin().WithClient(
		harbor.New(),
	).Run()

	if err != nil {
		fmt.Printf("plugin run err: %s, exit\n", err.Error())
		os.Exit(1)
	}
}
