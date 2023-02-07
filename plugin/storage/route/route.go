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

package route

import (
	"fmt"
	"path"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/storage"
	filestorev1alpha1 "github.com/katanomi/pkg/plugin/storage/capabilities/filestore/v1alpha1"
	"github.com/katanomi/pkg/plugin/storage/core/v1alpha1"
	v1alpha12 "github.com/katanomi/pkg/plugin/storage/route/core/v1alpha1"
	v1alpha13 "github.com/katanomi/pkg/plugin/storage/route/filestore/v1alpha1"
)

// GetPluginAPIPath returns storage plugin web service path
func GetPluginAPIPath(c client.Interface) string {
	return path.Join("/storage", c.Path())
}

// NewService new service from storage plugin client
func NewService(c client.Interface, filters ...restful.FilterFunction) (*restful.WebService, error) {
	routes := match(c)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no route for provider %s", c.Path())
	}

	group := &restful.WebService{}

	// adds standard prefix for plugins
	pluginAPIPath := GetPluginAPIPath(c)
	// group.Path(pluginAPIPath)
	for _, filter := range filters {
		group.Filter(filter)
	}

	for _, r := range routes {
		group.Path(path.Join(pluginAPIPath, r.GroupVersion().Identifier())).
			Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
		r.Register(group)
	}

	return group, nil
}

// match math route with plugin storage plugin client
func match(c client.Interface) []storage.VersionedRouter {
	routes := make([]storage.VersionedRouter, 0)

	if authCheck, ok := c.(v1alpha1.AuthChecker); ok {
		routes = append(routes, v1alpha12.NewAuthCheck(authCheck))
	}

	if filestore, ok := c.(filestorev1alpha1.FileStoreCapable); ok {
		routes = append(routes, v1alpha13.NewFileObject(filestore))
		routes = append(routes, v1alpha13.NewFileMeta(filestore))
	}
	return routes
}
