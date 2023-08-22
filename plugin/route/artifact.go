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
	"io"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/plugin/path"
	"knative.dev/pkg/logging"
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

type projectArtifactLister struct {
	impl client.ProjectArtifactLister
	tags []string
}

// NewProjectArtifactList create a list project artifacts route with plugin client
func NewProjectArtifactList(impl client.ProjectArtifactLister) Route {
	return &projectArtifactLister{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *projectArtifactLister) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integraion")
	ws.Route(
		ListOptionsDocs(
			ws.GET("/projects/{project:*}/artifacts").To(a.ListProjectArtifacts).
				// docs
				Doc("ListProjectArtifacts").Param(projectParam).
				Metadata(restfulspec.KeyOpenAPITags, a.tags).
				Returns(http.StatusOK, "OK", metav1alpha1.ArtifactList{}),
		),
	)
}

func (a *projectArtifactLister) ListProjectArtifacts(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	pathParams := metav1alpha1.ArtifactOptions{
		RepositoryOptions: metav1alpha1.RepositoryOptions{
			Project: path.Parameter(request, "project"),
		},
	}
	artifacts, err := a.impl.ListProjectArtifacts(request.Request.Context(), pathParams, option)
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

type projectArtifactUploader struct {
	impl client.ProjectArtifactUploader
	tags []string
}

// NewProjectArtifactUploader create an upload artifact route with plugin client
func NewProjectArtifactUploader(impl client.ProjectArtifactUploader) Route {
	return &projectArtifactUploader{
		tags: []string{"projects", "artifacts"},
		impl: impl,
	}
}

func (a *projectArtifactUploader) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integration")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.PUT("/projects/{project:*}/artifacts/{artifact:*}").To(a.UploadProjectArtifact).
			// docs
			Doc("UploadProjectArtifact").Param(projectParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.Artifact{}),
	)
}

// UploadProjectArtifact http handler for upload artifact
func (a *projectArtifactUploader) UploadProjectArtifact(request *restful.Request, response *restful.Response) {
	artifactOptions := metav1alpha1.ProjectArtifactOptions{
		Project:  request.PathParameter("project"),
		Artifact: request.PathParameter("artifact"),
	}
	artifactOptions.SubResourcesOptions = client.GetSubResourcesOptionsFromRequest(request)
	err := a.impl.UploadArtifact(request.Request.Context(), artifactOptions, request.Request.Body)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

type projectArtifactGetterDownloader struct {
	impl client.ProjectArtifactGetterDownloader
	tags []string
}

// NewArtifactGetterDownload create a get or download artifact route with plugin client
func NewArtifactGetterDownload(impl client.ProjectArtifactGetterDownloader) Route {
	return &projectArtifactGetterDownloader{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *projectArtifactGetterDownloader) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "project belong to integration")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	downloadParam := ws.QueryParameter("download", "download artifact")
	ws.Route(
		ws.GET("/projects/{project:*}/artifacts/{artifact:*}").To(a.GetOrDownloadProjectArtifact).
			// docs
			Doc("GetOrDownloadProjectArtifact").Param(projectParam).Param(artifactParam).Param(downloadParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Produces(restful.MIME_JSON, restful.MIME_OCTET).
			Returns(http.StatusOK, "OK", metav1alpha1.Artifact{}),
	)
}

// GetOrDownloadProjectArtifact http handler for get or download artifact detail
func (a *projectArtifactGetterDownloader) GetOrDownloadProjectArtifact(request *restful.Request, response *restful.Response) {
	logger := logging.FromContext(request.Request.Context())
	logger.Info("GetOrDownloadProjectArtifact")
	projectArtifactOptions := metav1alpha1.ProjectArtifactOptions{
		Project:  path.Parameter(request, "project"),
		Artifact: path.Parameter(request, "artifact"),
	}
	projectArtifactOptions.SubResourcesOptions = client.GetSubResourcesOptionsFromRequest(request)
	download := request.QueryParameter("download")
	if download != "" {
		readCloser, err := a.impl.DownloadProjectArtifact(request.Request.Context(), projectArtifactOptions)
		if err != nil {
			kerrors.HandleError(request, response, err)
			return
		}
		defer readCloser.Close()
		_, err = io.Copy(response.ResponseWriter, readCloser)
		if err != nil {
			kerrors.HandleError(request, response, err)
			return
		}
	} else {
		artifact, err := a.impl.GetProjectArtifact(request.Request.Context(), projectArtifactOptions)
		if err != nil {
			kerrors.HandleError(request, response, err)
			return
		}
		response.WriteHeaderAndEntity(http.StatusOK, artifact)
	}
}

type projectArtifactDeleter struct {
	impl client.ProjectArtifactDeleter
	tags []string
}

// NewProjectArtifactDeleter create a delete artifact route with plugin client
func NewProjectArtifactDeleter(impl client.ProjectArtifactDeleter) Route {
	return &projectArtifactDeleter{
		tags: []string{"projects", "repositories", "artifacts"},
		impl: impl,
	}
}

func (a *projectArtifactDeleter) Register(ws *restful.WebService) {
	projectParam := ws.PathParameter("project", "repository belong to integraion")
	repositoryParam := ws.PathParameter("repository", "artifact belong to repository")
	artifactParam := ws.PathParameter("artifact", "artifact name, maybe is version or tag")
	ws.Route(
		ws.DELETE("/projects/{project:*}/artifacts/{artifact:*}").To(a.DeleteArtifact).
			// docs
			Doc("DeleteArtifact").Param(projectParam).Param(repositoryParam).Param(artifactParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", nil),
	)
}

// DeleteArtifact http handler for delete artifact
func (a *projectArtifactDeleter) DeleteArtifact(request *restful.Request, response *restful.Response) {
	artifactOpts := metav1alpha1.ProjectArtifactOptions{
		Project:  path.Parameter(request, "project"),
		Artifact: path.Parameter(request, "artifact"),
	}
	artifactOpts.SubResourcesOptions = client.GetSubResourcesOptionsFromRequest(request)
	err := a.impl.DeleteProjectArtifact(request.Request.Context(), artifactOpts)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeader(http.StatusOK)
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
