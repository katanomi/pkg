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

package registry

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"knative.dev/pkg/logging"
)

// RegistrySchemeDetection detect registry scheme
type RegistrySchemeDetection interface {
	// DetectScheme detect registry scheme
	DetectScheme(ctx context.Context, registry, username, password string) (string, error)

	// DetectSchemeWithDefault detect registry scheme, if detect failed, return default scheme
	DetectSchemeWithDefault(ctx context.Context, registry, username, password, defaultScheme string) string
}

// DefaultRegistrySchemeDetection default scheme detection
type DefaultRegistrySchemeDetection struct {
	// Client is the resty client
	Client *resty.Client

	// Server should be accessed without verifying the TLS certificate.
	Insecure bool

	// Cache indicates whether to cache the scheme
	Cache       bool
	schemeCache sync.Map

	// httpClient for testing only
	httpClient *http.Client
}

// NewDefaultRegistrySchemeDetection create default registry scheme detection
func NewDefaultRegistrySchemeDetection(client *resty.Client, insecure, cache bool) *DefaultRegistrySchemeDetection {
	return &DefaultRegistrySchemeDetection{
		Client:   client,
		Insecure: insecure,
		Cache:    cache,
	}
}

// DefaultDetectImageRegistryScheme detect image registry scheme using default registry pinger
func DefaultDetectImageRegistryScheme(registryHost string, registryClient *http.Client, insecure bool) (string, error) {
	registryPinger := &DefaultRegistryPinger{
		Client:   registryClient,
		Insecure: insecure,
	}

	// verify the registry connection now to avoid future surprises
	registryURL, err := registryPinger.Ping(registryHost)
	if err != nil {
		return "", fmt.Errorf("failed to ping registry %s: %v", registryHost, err)
	}
	return registryURL.Scheme, nil
}

// DetectScheme detect registry scheme
func (d *DefaultRegistrySchemeDetection) DetectScheme(ctx context.Context, registry, username, password string) (string, error) {

	if strings.HasPrefix(registry, "http://") {
		return "http", nil
	}
	if strings.HasPrefix(registry, "https://") {
		return "https", nil
	}

	log := logging.FromContext(ctx).With("registry", registry)
	if d.Cache {
		if scheme, ok := d.schemeCache.Load(registry); ok {
			log.Debugw("get registry scheme from cache", "scheme", scheme)
			return scheme.(string), nil
		}
	}

	if d.Client == nil {
		return "", fmt.Errorf("registry client is nil")
	}

	httpClient := d.Client.GetClient()
	if username != "" && password != "" {
		httpClient = d.Client.SetBasicAuth(username, password).GetClient()
	}
	if d.httpClient != nil {
		// for testing only
		httpClient = d.httpClient
	}

	scheme, err := DefaultDetectImageRegistryScheme(registry, httpClient, d.Insecure)
	if err != nil {
		log.Errorw("failed to detect registry scheme", "error", err)
		return "", err
	}
	log.Infow("detect registry scheme", "scheme", scheme)

	if d.Cache {
		d.schemeCache.Store(registry, scheme)
	}
	return scheme, nil
}

// DetectSchemeWithDefault detect registry scheme, if detect failed, return default scheme
func (d *DefaultRegistrySchemeDetection) DetectSchemeWithDefault(ctx context.Context, registry, username, password, defaultScheme string) string {
	scheme, err := d.DetectScheme(ctx, registry, username, password)
	if err != nil {
		return defaultScheme
	}
	return scheme
}
