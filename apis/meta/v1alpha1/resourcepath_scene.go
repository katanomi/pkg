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

package v1alpha1

// ResourcePathScene resource path scene for formatting resource path
type ResourcePathScene string

const (
	// ResourcePathSceneAPI resource path scene for api call
	ResourcePathSceneAPI ResourcePathScene = "api"
	// ResourcePathSceneWebConsole resource path scene for web page
	ResourcePathSceneWebConsole ResourcePathScene = "web-console"
	// ResourcePathSceneHttpClone resource path scene for http clone url for git
	ResourcePathSceneHttpClone ResourcePathScene = "http-clone"
)
