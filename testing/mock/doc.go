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

// Package mock contains useful functionality for testing
package mock

//go:generate mockgen -package=client -destination=./sigs.k8s.io/controller-runtime/pkg/client/client.go sigs.k8s.io/controller-runtime/pkg/client Client
//go:generate mockgen -package=manager -destination=./sigs.k8s.io/controller-runtime/pkg/manager/manager.go sigs.k8s.io/controller-runtime/pkg/manager Manager
//go:generate mockgen -package=apis -destination=./knative.dev/pkg/apis/condition_manager.go knative.dev/pkg/apis ConditionManager
//go:generate mockgen -package=kubernetes -destination=./k8s.io/client-go/kubernetes/clientset.go k8s.io/client-go/kubernetes Interface
//go:generate mockgen -package=sharedmain -destination=./github.com/katanomi/pkg/watcher.go github.com/katanomi/pkg/sharedmain DefaultingWatcherWithOnChange
