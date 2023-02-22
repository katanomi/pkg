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

package secret

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
)

// NewResourcePathFormat construct a ResourcePathFormat from json string
func NewResourcePathFormat(pathJson, subPathJson string) *ResourcePathFormat {
	r := &ResourcePathFormat{}
	_ = json.Unmarshal([]byte(pathJson), &r.pathData)
	_ = json.Unmarshal([]byte(subPathJson), &r.subPathData)
	return r
}

// ResourcePathFormat manage the resource path format
type ResourcePathFormat struct {
	pathData    map[metav1alpha1.ResourcePathScene]string
	subPathData map[metav1alpha1.ResourcePathScene]string
}

// GetFmt get the format string for a special scene
func (p *ResourcePathFormat) getPathFmt(scene metav1alpha1.ResourcePathScene) string {
	if p.pathData == nil {
		return ""
	}
	return p.pathData[scene]
}

// GetFmt get the format string for a special scene
func (p *ResourcePathFormat) getSubPathFmt(scene metav1alpha1.ResourcePathScene) string {
	if p.subPathData == nil {
		return ""
	}
	return p.subPathData[scene]
}

func (p *ResourcePathFormat) splitScope(scope string) []string {
	_segments := strings.Split(scope, "/")
	var segments []string
	for _, item := range _segments {
		if item == "" {
			continue
		}
		segments = append(segments, item)
	}
	return segments
}

func (p *ResourcePathFormat) formatResourcePath(scene metav1alpha1.ResourcePathScene, segments []string) string {
	pathFmt := p.getPathFmt(scene)
	if pathFmt == "" {
		pathFmt = "%s"
	}
	pathFmt = "/" + strings.Trim(pathFmt, "/") + "/"
	return fmt.Sprintf(pathFmt, segments[0])
}

func (p *ResourcePathFormat) formatSubResourcePath(scene metav1alpha1.ResourcePathScene, segments []string) string {
	subPathFmt := p.getSubPathFmt(scene)
	if subPathFmt == "" {
		subPathFmt = "%s/%s"
	}
	subPathFmt = "/" + strings.Trim(subPathFmt, "/") + "/"
	resource := segments[0]
	subRes := strings.Join(segments[1:], "/")
	return fmt.Sprintf(subPathFmt, resource, subRes)
}

// FormatPathByScene get the formatted string of special scene
func (p *ResourcePathFormat) FormatPathByScene(scene metav1alpha1.ResourcePathScene, scope string) string {
	segments := p.splitScope(scope)
	if len(segments) == 0 {
		return "/"
	}

	if len(segments) == 1 {
		return p.formatResourcePath(scene, segments)
	}

	return p.formatSubResourcePath(scene, segments)
}

func (p *ResourcePathFormat) getMapKeys(m map[metav1alpha1.ResourcePathScene]string) []metav1alpha1.ResourcePathScene {
	var sceneList []metav1alpha1.ResourcePathScene
	for scene := range m {
		sceneList = append(sceneList, scene)
	}
	sort.Slice(sceneList, func(i, j int) bool {
		return sceneList[i] < sceneList[j]
	})
	return sceneList
}

// FormatPathAllScene get the formatted strings of all scenes
func (p *ResourcePathFormat) FormatPathAllScene(scope string) (list []string) {
	segments := p.splitScope(scope)
	if len(segments) == 0 {
		return []string{"/"}
	}
	defer func() {
		if len(list) == 0 {
			list = append(list, scope)
		}
	}()
	if len(segments) == 1 {
		sceneList := p.getMapKeys(p.pathData)
		for _, scene := range sceneList {
			list = append(list, p.formatResourcePath(scene, segments))
		}
		return list
	}

	sceneList := p.getMapKeys(p.subPathData)
	for _, scene := range sceneList {
		list = append(list, p.formatSubResourcePath(scene, segments))
	}
	return list
}
