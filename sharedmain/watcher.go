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

package sharedmain

import (
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
)

// DefaultingWatcherWithOnChange is a configmap.DefaultingWatcher that also has an OnChange callback.
type DefaultingWatcherWithOnChange interface {
	// DefaultingWatcher is similar to Watcher, but if a ConfigMap is absent, then a code provided
	// default will be used.
	configmap.DefaultingWatcher

	// OnChange invokes the callbacks of all observers of the given ConfigMap.
	OnChange(*corev1.ConfigMap)
}
