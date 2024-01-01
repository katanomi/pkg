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

package conformance

import (
	"fmt"
	"testing"
)

func TestAbc(t *testing.T) {
	a := &Node{Label: newTestCaseLabel("a")}
	b := &Node{Label: newTestCaseLabel("b")}
	c := &Node{Label: newTestCaseLabel("c")}
	d := &Node{Label: newTestCaseLabel("d")}
	d.LinkParentNode(c)
	c.LinkParentNode(a)
	b.LinkParentNode(a)

	fmt.Println(a.FullPathLabels(), d.IdentifyLabels())
	e := a.Clone()
	fmt.Println(e)
}
