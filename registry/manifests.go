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

package registry

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/mod"
	"github.com/regclient/regclient/types/docker/schema2"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/ref"
)

type ManifestClient struct {
	options         []regclient.Opt
	schemeDetection RegistrySchemeDetection
	insecure        bool
}

func NewManifestClient(options ...regclient.Opt) *ManifestClient {
	return &ManifestClient{
		options:         options,
		insecure:        true,
		schemeDetection: NewDefaultRegistrySchemeDetection(resty.New(), true, true),
	}
}

// Insecure access registry without verifying the TLS certificate
func (c *ManifestClient) Insecure(value bool) *ManifestClient {
	if c.insecure != value {
		c.insecure = value
		c.schemeDetection = NewDefaultRegistrySchemeDetection(resty.New(), value, true)
	}

	return c
}

// GetAnnotations get annotations from a reference image
func (c *ManifestClient) GetAnnotations(ctx context.Context, reference string) (map[string]string, error) {
	r, err := ref.New(reference)
	if err != nil {
		return nil, err
	}

	regClient := c.newRegClient(ctx, r)
	defer regClient.Close(ctx, r)

	m, err := regClient.ManifestGet(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("get manifest error: %s", err.Error())
	}

	if anno, ok := m.(manifest.Annotator); ok {
		return anno.GetAnnotations()
	}

	return nil, fmt.Errorf("manifest can not access annotation")
}

// PutEmptyIndex create an empty manifest list with annotations
func (c *ManifestClient) PutEmptyIndex(ctx context.Context, reference string, annotations map[string]string) error {
	r, err := ref.New(reference)
	if err != nil {
		return err
	}

	regClient := c.newRegClient(ctx, r)
	defer regClient.Close(ctx, r)

	options := []manifest.Opts{
		manifest.WithOrig(schema2.ManifestList{
			Versioned:   schema2.ManifestListSchemaVersion,
			Annotations: annotations,
		}),
	}

	m, err := manifest.New(options...)
	if err != nil {
		return err
	}

	return regClient.ManifestPut(ctx, r, m)
}

// SetAnnotation append annotation to  a reference image
// annotation key will be deleted with empty value
func (c *ManifestClient) SetAnnotation(ctx context.Context, reference string, annotations map[string]string) error {
	r, err := ref.New(reference)
	if err != nil {
		return err
	}

	regClient := c.newRegClient(ctx, r)
	defer regClient.Close(ctx, r)

	if r.Tag == "" {
		return fmt.Errorf("cannot replace an image digest, must include a tag")
	}

	modOptions := make([]mod.Opts, 0)
	for key, value := range annotations {
		modOptions = append(modOptions, mod.WithAnnotation(key, value))
	}

	output, err := mod.Apply(ctx, regClient, r, modOptions...)
	if err != nil {
		return fmt.Errorf("apply annotation error: %s", err.Error())
	}

	err = regClient.ImageCopy(ctx, output, r)
	if err != nil {
		return fmt.Errorf("failed copying image to new name: %w", err)
	}

	return nil
}

// overrideHostTLS tls can only be overridden by options as RegClient does not provide a set host method
func (c *ManifestClient) overrideHostTLS(ctx context.Context, reference string) regclient.Opt {
	host := config.HostNewName(reference)

	scheme, _ := c.schemeDetection.DetectScheme(ctx, host.Name)

	if scheme == HTTP {
		host.TLS = config.TLSDisabled
	} else if c.insecure {
		host.TLS = config.TLSInsecure
	} else {
		host.TLS = config.TLSEnabled
	}

	return regclient.WithConfigHost(*host)
}

func (c *ManifestClient) newRegClient(ctx context.Context, ref ref.Ref) *regclient.RegClient {
	options := append([]regclient.Opt{
		regclient.WithDockerCreds(),
		c.overrideHostTLS(ctx, ref.Reference),
	}, c.options...)

	return regclient.New(options...)
}
