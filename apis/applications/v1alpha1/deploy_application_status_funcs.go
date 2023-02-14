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
	corev1 "k8s.io/api/core/v1"
)

// IsEmpty returns true if there is not enought information
// regarding the application deployment
func (app *DeployApplicationResults) IsEmpty() bool {
	return app == nil || isAppRefEmpty(app.ApplicationRef) || (isDeployApplicationStatusEmpty(app.After) && isDeployApplicationStatusEmpty(app.Before))
}

func isAppRefEmpty(ref *corev1.ObjectReference) bool {
	return ref == nil || (ref.Kind == "" && ref.Name == "" && ref.Namespace == "")
}

func isDeployApplicationStatusEmpty(items []DeployApplicationStatus) bool {
	for _, item := range items {
		if item.Name != "" || item.Version != "" || item.Status != "" {
			return false
		}
	}
	return true
}
