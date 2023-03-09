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

package workers

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// JobRunnable provides runnable methods to managed by worker
type JobRunnable interface {
	Setup(ctx context.Context, kclient client.Client, restClient *resty.Client) error
	// JobName is identical to config key name from configManager
	JobName() string
	// RunFunc returns the cron callback func to be invoked
	RunFunc(ctx context.Context) func()
}

// cronJob handles cron job object at runtime
type cronJob struct {
	// Name for cron job name identical to ConfigManager data key
	name string

	// FuncJob for job func
	funcJob cron.FuncJob

	// entryID for cron EntryID
	entryID cron.EntryID
}
