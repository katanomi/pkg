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

package names

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

// GetGenerateName used to generate child object names with a "-" suffix
func GetGenerateName(object metav1.Object) string {
	name := object.GetName()
	if name == "" {
		name = object.GetGenerateName()
	}

	return fmt.Sprintf("%s-", name)
}

const (
	// TODO: make this flexible for non-core resources with alternate naming rules.
	maxNameLength          = 63
	randomLength           = 5
	MaxGeneratedNameLength = maxNameLength - randomLength
)

// copied from https://github.com/kubernetes/kubernetes/blob/c9fb3c8a1b3f407a5e84562843780aa3047d7d06/staging/src/k8s.io/apiserver/pkg/storage/names/generate.go#L49
// temporarily

// Generates a name with a random suffix
func GenerateName(base string) string {
	if len(base) > MaxGeneratedNameLength {
		base = base[:MaxGeneratedNameLength]
	}
	return fmt.Sprintf("%s%s", base, utilrand.String(randomLength))
}
