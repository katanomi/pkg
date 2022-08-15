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

package manager

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"sigs.k8s.io/controller-runtime/pkg/leaderelection"
	"sigs.k8s.io/controller-runtime/pkg/recorder"
)

// ResourceLockFunc resouce lock function
// Ref: https://github.com/kubernetes-sigs/controller-runtime/blob/1638a6a9b82dc1e0046c7a1006f12dacd9475f54/pkg/leaderelection/leader_election.go#L54
type ResourceLockFunc func(*rest.Config, recorder.Provider, leaderelection.Options) (resourcelock.Interface, error)

// ResourceLockSetter sets a new resource lock function
// this will make possible to customize resource lock without
// having to force it into dependencies
type ResourceLockSetter interface {
	SetNewResourceLock(lock ResourceLockFunc)
}
