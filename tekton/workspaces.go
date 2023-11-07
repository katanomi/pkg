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

package tekton

import (
	"context"
	"fmt"
	"strings"

	kclient "github.com/katanomi/pkg/client"
	pkgerr "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/storage/storageclass"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetWorkspaceBindings returns workspace binding from pipeineRun
func GetWorkspaceBindings(ctx context.Context, pr *v1beta1.PipelineRun) ([]v1beta1.WorkspaceBinding, error) {
	log := logging.FromContext(ctx)
	clt := kclient.Client(ctx)

	if clt == nil {
		return nil, fmt.Errorf("nil k8s client from context")
	}

	var (
		workspaceBindings = make([]v1beta1.WorkspaceBinding, 0, len(pr.Spec.Workspaces))
		templateBindings  = make([]v1beta1.WorkspaceBinding, 0)
	)

	for _, ws := range pr.Spec.Workspaces {
		if ws.VolumeClaimTemplate == nil {
			workspaceBindings = append(workspaceBindings, ws)
			continue
		}
		templateBindings = append(templateBindings, ws)
	}

	matchLabels := &client.MatchingLabels{
		pipeline.PipelineRunLabelKey: pr.Name,
	}

	taskRunList := &v1beta1.TaskRunList{}
	if err := clt.List(ctx, taskRunList, matchLabels, client.InNamespace(pr.Namespace)); err != nil {
		log.Errorw("list taskrun error", "error", err)
		return nil, err
	}

	taskWorkspaceBindingMap := getWorkspacePipelineTaskBindingMap(pr)

	log.Debugw("invoke getWorkspacePipelineTaskBindingMap in GetWorkspaceBindings", "map", taskWorkspaceBindingMap)

	for _, templateWorkspace := range templateBindings {
		// for each template workspace, try to find the binding in taskrun
		binding := findTaskRunWorkspaceBindingFromTaskRuns(templateWorkspace.Name, taskRunList.Items, taskWorkspaceBindingMap)
		if binding != nil {
			templateWorkspace.PersistentVolumeClaim = binding.PersistentVolumeClaim
		}
		workspaceBindings = append(workspaceBindings, templateWorkspace)
	}
	log.Debugw("GetWorkspaceBindings return binding", "workspaceBindings", workspaceBindings)

	return workspaceBindings, nil
}

// CheckWorkspaceBindings check if workspace bindings are valid
func CheckWorkspaceBindings(ctx context.Context, pr *v1beta1.PipelineRun) error {

	if pr == nil {
		return pkgerr.ErrNilPointer
	}

	workspaces := pr.Spec.Workspaces
	if len(workspaces) == 0 {
		return nil
	}

	var emptySCVolumeClaimTemplateWSNames []string

	for _, ws := range workspaces {
		if ws.VolumeClaimTemplate != nil {
			if ws.VolumeClaimTemplate.Spec.StorageClassName == nil ||
				*ws.VolumeClaimTemplate.Spec.StorageClassName == "" {
				emptySCVolumeClaimTemplateWSNames = append(emptySCVolumeClaimTemplateWSNames, ws.Name)
			}
		}
	}

	if len(emptySCVolumeClaimTemplateWSNames) == 0 {
		return nil
	}

	cli := kclient.Client(ctx)
	if cli == nil {
		return fmt.Errorf("nil k8s client from context")
	}
	defaultStorageClass := storageclass.GetDefaultStorageClass(ctx, cli)
	if defaultStorageClass == nil {
		return pkgerr.NewDefaultStorageClassNotFound(
			fmt.Sprintf(
				"no default storageclass found, workspaces use volumeClaimTemplate but no storageclass name: %s",
				strings.Join(emptySCVolumeClaimTemplateWSNames, ", "),
			),
		)
	}
	return nil
}

func findTaskRunWorkspaceBindingFromTaskRuns(pipelineWorkspaceName string,
	taskRuns []v1beta1.TaskRun,
	wsPipelineTaskBindingMap map[string]map[string]string) *v1beta1.WorkspaceBinding {
	for _, taskRun := range taskRuns {
		if len(taskRun.Spec.Workspaces) == 0 {
			continue
		}
		taskRunLabels := taskRun.GetLabels()
		if taskRunLabels == nil {
			continue
		}

		pipelineTaskName, exists := taskRunLabels[pipeline.PipelineTaskLabelKey]
		if !exists || pipelineTaskName == "" {
			continue
		}

		taskWorkspaceBinding, taskExists := wsPipelineTaskBindingMap[pipelineTaskName]
		if !taskExists {
			continue
		}

		taskRunWorkspaceName, wsExists := taskWorkspaceBinding[pipelineWorkspaceName]
		if !wsExists || taskRunWorkspaceName == "" {
			continue
		}

		for _, trWS := range taskRun.Spec.Workspaces {
			if trWS.Name == taskRunWorkspaceName {
				return &trWS
			}
		}
	}
	return nil
}

// getWorkspacePipelineTaskBindingMap return workspaces binding map on pipeline tasks mapped by task name
// ret: taskName => bindingMap
// bindingMap: pipelineWorkspaceName => taskRunWorkspaceName
func getWorkspacePipelineTaskBindingMap(pr *v1beta1.PipelineRun) map[string]map[string]string {
	pipelineSpec := GetPipelineSpec(pr)
	if pipelineSpec == nil {
		return map[string]map[string]string{}
	}
	ret := make(map[string]map[string]string)
	for _, task := range pr.Spec.PipelineSpec.Tasks {
		mapItem := make(map[string]string)
		for _, wsBinding := range task.Workspaces {
			mapItem[wsBinding.Workspace] = wsBinding.Name
		}
		ret[task.Name] = mapItem
	}
	return ret
}
