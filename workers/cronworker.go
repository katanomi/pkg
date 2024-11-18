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

	metav1alpha1 "github.com/AlaudaDevops/pkg/apis/meta/v1alpha1"
	"github.com/AlaudaDevops/pkg/config"
	"github.com/AlaudaDevops/pkg/errors"
	"github.com/AlaudaDevops/pkg/restclient"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// ConfigWatcherFunc is the default watch func
var ConfigWatcherFunc = func(cw *CronWorker) func(c *config.Config) {
	return func(c *config.Config) {
		for i, job := range cw.jobs {
			spec := metav1alpha1.DataMap(c.Data).MustStringVal(job.name, "0 0 * * *")
			newEntryID, err := cw.cron.AddJob(spec, job.funcJob)
			if err != nil {
				cw.Errorw("config watcher update cron job error", "err", err)
				return
			}
			if job.entryID > 0 {
				cw.cron.Remove(job.entryID)
			}
			cw.jobs[i].entryID = newEntryID
			cw.Debugf("ConfigWatcherFunc fallback: set job %s cron spec to %s", job.name, spec)
		}
	}
}

// CronWorker maintains a list of runners and dedicated to cron watched by configManager
type CronWorker struct {
	Runners []JobRunnable

	jobs []*cronJob

	*zap.SugaredLogger
	cron    *cron.Cron
	watcher config.Watcher
}

// NeedLeaderElection indicates cron worker
func (cw *CronWorker) NeedLeaderElection() bool {
	return true
}

// Start starts cron and waits for context cancellation
func (cw *CronWorker) Start(ctx context.Context) error {
	cw.cron.Start()
	<-ctx.Done()
	cw.cron.Stop()
	return nil
}

func (cw *CronWorker) Name() string {
	return "cron-worker"
}

func (cw *CronWorker) Setup(ctx context.Context, manager manager.Manager, logger *zap.SugaredLogger) error {
	restClient := restclient.RESTClient(ctx)

	if restClient == nil {
		return errors.ErrNilPointer
	}

	cw.cron = cron.New()
	cw.SugaredLogger = logger.With("component", cw.Name())
	ctx = logging.WithLogger(ctx, logger)

	for _, j := range cw.Runners {
		err := j.Setup(ctx, manager.GetClient(), restClient)
		if err != nil {
			return err
		}
		cw.jobs = append(cw.jobs, &cronJob{
			name:    j.JobName(),
			funcJob: j.RunFunc(ctx),
		})
	}

	if kMgr := config.KatanomiConfigManager(ctx); kMgr != nil {
		kMgr.AddWatcher(config.NewConfigWatcher(ConfigWatcherFunc(cw)))
	}

	return manager.Add(cw)
}

func (cw *CronWorker) CheckSetup(ctx context.Context, manager manager.Manager, logger *zap.SugaredLogger) error {
	return nil
}
