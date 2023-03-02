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

// Package cleanup for clean pvc mount for taskrun
package cleanup

import (
	"context"

	pkgmetav1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// PVCFinalizersForPod Finalizers in pvc about pod
	PVCFinalizersForPod = "kubernetes.io/pvc-protection"
	// ConditionCleanup is triggered when doing GC
	ConditionCleanup = "Cleanup"
)

type ShouldSkipCleanup func(ctx context.Context, workspace tekton.WorkspaceBinding) bool

func IsVolumeCreatedManaualy(workspaces []tekton.WorkspaceBinding, v tekton.WorkspaceBinding) bool {
	for _, ws := range workspaces {
		if ws.PersistentVolumeClaim != nil && v.PersistentVolumeClaim != nil &&
			ws.PersistentVolumeClaim.ClaimName == v.PersistentVolumeClaim.ClaimName {
			return true
		}
	}
	return false
}

// CleanTaskRunsPVC clean taskrun workspace pvc that dynamic create
// Params:
// taskLabels for filt taskrun with label, maybe is empty
// shouldSkip for skip clean
// Results:
// succeeded, fails indicate clean pvc successful and failed list
func CleanTaskRunsPVC(ctx context.Context,
	taskLables client.MatchingLabels,
	shouldSkip ShouldSkipCleanup,
) (succeeded []types.NamespacedName, fails map[string]error, err error) {
	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)

	taskRunList := &tekton.TaskRunList{}
	if err = clt.List(ctx, taskRunList, taskLables); err != nil {
		log.Errorw("list taskrun error", "error", err)
		return nil, nil, err
	}

	succeeded = make([]types.NamespacedName, 0)
	fails = make(map[string]error)
	errs := []error{}
	for _, taskRun := range taskRunList.Items {
		if err := cleanTaskRunPVC(ctx, &taskRun, shouldSkip, &succeeded, fails); err != nil {
			errs = append(errs, err)
		}
	}

	return succeeded, fails, utilerrors.NewAggregate(errs)
}

func cleanTaskRunPVC(
	ctx context.Context,
	taskRun *tekton.TaskRun,
	shouldSkip ShouldSkipCleanup,
	succeeded *[]types.NamespacedName,
	fails map[string]error,
) error {
	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)

	errs := []error{}
	log.Debugf("clean taskrun pvc", "taskrun", client.ObjectKey{Namespace: taskRun.Namespace, Name: taskRun.Name})
	for _, workspace := range taskRun.Spec.Workspaces {

		if shouldSkip(ctx, workspace) {
			log.Debugw("workspace is skip, no need to clean", "workspace", workspace)
			continue
		}

		pvcKey := client.ObjectKey{Namespace: taskRun.Namespace, Name: workspace.PersistentVolumeClaim.ClaimName}
		if pvcCleaned(*succeeded, pvcKey) {
			continue
		}

		log := log.With("pvc", pvcKey)
		ctx = logging.WithLogger(ctx, log)
		var pvc corev1.PersistentVolumeClaim
		if err := clt.Get(ctx, pvcKey, &pvc); err != nil {
			if !errors.IsNotFound(err) {
				log.Errorw("getting taskRun related pvc failed", "error", err)
				fails[pvcKey.String()] = err
				errs = append(errs, err)
			}
			continue
		}

		err := deletePVC(ctx, pvc)
		if err != nil {
			fails[pvcKey.String()] = err
			errs = append(errs, err)
		} else {
			*succeeded = append(*succeeded, pkgmetav1.GetNamespacedNameFromObject(&pvc))
		}
	}

	return utilerrors.NewAggregate(errs)
}

func deletePVC(ctx context.Context, pvc corev1.PersistentVolumeClaim) error {
	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)
	log.Info("delete taskrun pvc")

	retryLimit := 3

	for i := 1; i <= retryLimit; i++ {

		if err := clt.Delete(ctx, &pvc); err != nil {
			log.Errorw("delete taskRun related pvc failed", "error", err, "retryNum", i)
			if i < retryLimit {
				continue
			}
			return err
		}

		patch := client.MergeFrom(pvc.DeepCopy())
		length := len(pvc.Finalizers) - 1
		for index := length; index >= 0; index-- {
			// for protect pvc add finalizer kubernetes.io/pvc-protection by k8s
			// if not delete it, pvc will always terminating until pod be deleted
			// but k8s doesn't auto delete pod is completed
			if pvc.Finalizers[index] == PVCFinalizersForPod {
				pvc.Finalizers = append(pvc.Finalizers[:index], pvc.Finalizers[index+1:]...)
			}
		}

		if err := clt.Patch(ctx, &pvc, patch); err != nil {
			log.Errorw("delete taskRun related pvc failed when patching pvc finalizer", "error", err, "retryNum", i)
			if i < retryLimit {
				continue
			}
			return err
		}
		break
	}
	return nil
}

// pvcCleaned return true if pvc has been cleaned
func pvcCleaned(succeeded []types.NamespacedName, pvcKey client.ObjectKey) bool {
	for _, alreadyCleanupPVC := range succeeded {
		if alreadyCleanupPVC.String() == pvcKey.String() {
			return true
		}
	}
	return false
}
