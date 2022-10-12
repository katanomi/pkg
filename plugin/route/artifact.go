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

package route

import (
	"net/http"

	"github.com/katanomi/pkg/plugin/path"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
)

type artifactList struct {
	impl client.ArtifactLister
	tags []string
}

// NewArtifactList create a list artifact route with plugin client
func NewArtifactList(impl client.ArtifactLister) Route {
	return &artifactList{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *artifactList) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/repositories/{repository:*}/artifacts").To(a.ListArtifacts).
				// docs
				Doc("ListArtifacts").Param(projectParam).Param(repositoryParam).
				Metadata(restfulspec.KeyOpenAPITags, a.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.ArtifactList{}),
		),
	)
}

func (a *artifactList) ListArtifacts(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
		Repository: path.Parameter(request, "repository"),
	}
	artifacts, err := a.impl.ListArtifacts(request.Request.Context(), pathParams, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, artifacts)
}

type artifactGetter struct {
	impl client.ArtifactGetter
	tags []string
}

// NewArtifactGet create a get artifact route with plugin client
func NewArtifactGet(impl client.ArtifactGetter) Route {
	return &artifactGetter{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *artifactGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.GET("/projects/{project:*}/repositories/{repository:*}/artifacts/{artifact}").To(a.GetArtifact).
			// docs
			Doc("GetArtifact").Param(projectParam).Param(repositoryParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.Artifact{}),
	)
}

// GetArtifact http handler for get artifact detail
func (a *artifactGetter) GetArtifact(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
		Repository: path.Parameter(request, "repository"),
		Artifact:   path.Parameter(request, "artifact"),
	}
	artifact, err := a.impl.GetArtifact(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, artifact)
}

type artifactDeleter struct {
	impl client.ArtifactDeleter
	tags []string
}

// NewArtifactDelete create a delete artifact route with plugin client
func NewArtifactDelete(impl client.ArtifactDeleter) Route {
	return &artifactDeleter{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *artifactDeleter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.DELETE("/projects/{project:*}/repositories/{repository:*}/artifacts/{artifact}").To(a.DeleteArtifact).
			// docs
			Doc("DeleteArtifact").Param(projectParam).Param(repositoryParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// DeleteArtifact http handler for delete artifact
func (a *artifactDeleter) DeleteArtifact(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
		Repository: path.Parameter(request, "repository"),
		Artifact:   path.Parameter(request, "artifact"),
	}
	err := a.impl.DeleteArtifact(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type artifactTagDeleter struct {
	impl client.ArtifactTagDeleter
	tags []string
}

// NewArtifactTagDelete create a delete artifact tag route with plugin client
func NewArtifactTagDelete(impl client.ArtifactTagDeleter) Route {
	return &artifactTagDeleter{
		tags: []string{"projects", "repositories", "artifacts", "tag"},
		impl: impl,
	}
}

func (a *artifactTagDeleter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	tagParam := ws.PathParameter("tag", "the name of the tag")

	ws.Route(
		ws.DELETE("/projects/{project:*}/repositories/{repository:*}/artifacts/{artifact}/tags/{tag}").To(a.DeleteArtifactTag).
			// docs
			Doc("DeleteArtifactTag").Param(projectParam).Param(repositoryParam).Param(artifactParam).Param(tagParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// DeleteArtifact http handler for delete artifact
func (a *artifactTagDeleter) DeleteArtifactTag(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.ArtifactTagOptions{
		ArtifactOptions: metav1alpha1.ArtifactOptions{
			RepositoryOptions: metav1alpha1.RepositoryOptions{
				Project: path.Parameter(request, "project"),
			},
			Repository: path.Parameter(request, "repository"),
			Artifact:   path.Parameter(request, "artifact"),
		},
		Tag: path.Parameter(request, "tag"),
	}

	err := a.impl.DeleteArtifactTag(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type scanImage struct {
	impl client.ScanImage
	tags []string
}

// NewScanImage create a scan image route with plugin client
func NewScanImage(impl client.ScanImage) Route {
	return &scanImage{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (s *scanImage) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.POST("/projects/{project:*}/repositories/{repository:*}/artifacts/{artifact}/scan").To(s.ScanImage).
			// docs
			Doc("ScanImage").Param(projectParam).Param(repositoryParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, s.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// ScanImage http handler for scan image
func (s *scanImage) ScanImage(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
		Repository: path.Parameter(request, "repository"),
		Artifact:   path.Parameter(request, "artifact"),
	}
	err := s.impl.ScanImage(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type imageConfigGetter struct {
	impl client.ImageConfigGetter
	tags []string
}

// NewImageConifgGetter create a get image config route with plugin client
func NewImageConifgGetter(impl client.ImageConfigGetter) Route {
	return &imageConfigGetter{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (i *imageConfigGetter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.GET("/projects/{project:*}/repositories/{repository:*}/artifacts/{artifact}/config").To(i.GetImageConfig).
			// docs
			Doc("GetArtifact").Param(projectParam).Param(repositoryParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, i.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.ImageConfig{}),
	)
}

// GetImageConfig http handler for get image config
func (i *imageConfigGetter) GetImageConfig(request *restful.Request, response *restful.Response) {
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
		Repository: path.Parameter(request, "repository"),
		Artifact:   path.Parameter(request, "artifact"),
	}
	config, err := i.impl.GetImageConfig(request.Request.Context(), pathParams)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, config)
}
