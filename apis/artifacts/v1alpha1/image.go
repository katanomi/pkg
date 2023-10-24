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

import (
	"context"
	"fmt"
	"net/http"

	"github.com/katanomi/pkg/client"
	pclient "github.com/katanomi/pkg/plugin/client"
	"github.com/katanomi/pkg/restclient"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"knative.dev/pkg/logging"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

// PushEmptyImage an empty image will be created and pushed.
// The function will read the credentials from the context. If the credentials
// are empty, only the credentials will not be used for push.
func PushEmptyImage(ctx context.Context, uri URI) error {
	log := logging.FromContext(ctx).Named("PushImage").With("ref", uri.String())
	if uri.Host == "" || uri.Path == "" || uri.Tag == "" {
		err := fmt.Errorf("repository host, path and tag must be set")
		return err
	}

	// create empty layers.
	memStore := memory.New()
	opts := oras.PackManifestOptions{
		Layers: []ocispec.Descriptor{},
	}

	manifestDescriptor, err := oras.PackManifest(ctx, memStore, oras.PackManifestVersion1_1_RC4, OCIEmptyArtifactType, opts)
	if err != nil {
		log.Errorw("failed to pack manifest with memory", "err", err)
		return err
	}
	log.Debugw("pack manifest descriptor", "descriptor", manifestDescriptor)

	if err = memStore.Tag(ctx, manifestDescriptor, uri.Tag); err != nil {
		log.Errorw("failed to create tag with memory", "err", err)
		return err
	}

	repo, err := remote.NewRepository(uri.Repository())
	if err != nil {
		log.Errorw("failed to new repository", "err", err)
		return err
	}

	credential, err := extraAuthFromContext(ctx)
	if err != nil {
		log.Warnw("failed to extra auth, will try to push without credentials", "err", err)
	}

	var httpCli *http.Client
	restyClient := restclient.RESTClient(ctx)
	if restyClient != nil {
		httpCli = restyClient.GetClient()
	} else {
		httpCli = client.NewHTTPClient()
	}

	// need to make sure ignore authentication is set.
	// Warn: may pollute the global client.
	client.InsecureSkipVerifyOption(httpCli)

	repo.Client = &auth.Client{
		Client:     httpCli,
		Cache:      auth.DefaultCache,
		Credential: auth.StaticCredential(uri.Host, credential),
	}

	// Copy from the memory store to the remote repository
	_, err = oras.Copy(ctx, memStore, uri.Tag, repo, uri.Tag, oras.DefaultCopyOptions)
	if err != nil {
		return err
	}
	return nil
}

func extraAuthFromContext(ctx context.Context) (auth.Credential, error) {
	authInfo := pclient.ExtractAuth(ctx)
	if authInfo == nil {
		return auth.EmptyCredential, nil
	}

	if !authInfo.IsBasic() {
		return auth.EmptyCredential, fmt.Errorf("only support basic auth, current auth type: %s", authInfo.Type)
	}

	username, password, err := authInfo.GetBasicInfo()
	if err != nil {
		return auth.EmptyCredential, err
	}

	return auth.Credential{
		Username: username,
		Password: password,
	}, nil
}
