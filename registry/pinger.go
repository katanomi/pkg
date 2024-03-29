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
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/katanomi/pkg/networkutils"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"knative.dev/pkg/logging"
)

// Ref: https://github.com/openshift/oc/blob/232e65f43c768fbe4f6a0438fc8671c69182b49a/pkg/cli/admin/prune/imageprune/helper.go

// RegistryPinger performs a health check against a registry.
type RegistryPinger interface {
	// Ping performs a health check against registry. It returns registry url qualified with schema unless an
	// error occurs.
	Ping(ctx context.Context, registry string) (*url.URL, error)
}

// DefaultRegistryPinger implements RegistryPinger.
type DefaultRegistryPinger struct {
	Client   *http.Client
	Insecure bool
}

// Ping verifies that the integrated registry is ready, determines its transport protocol and returns its url
// or error.
func (drp *DefaultRegistryPinger) Ping(ctx context.Context, registry string) (*url.URL, error) {
	var (
		registryURL *url.URL
		err         error
	)

pathLoop:
	// first try the new default / path, then fall-back to the obsolete /healthz endpoint
	for _, path := range []string{"/", "/healthz"} {
		registryURL, err = TryProtocolsWithRegistryURL(ctx, registry, drp.Insecure, func(u url.URL) error {
			u.Path = path
			req, err := http.NewRequest(http.MethodGet, u.String(), nil)
			if err != nil {
				return err
			}
			req = req.WithContext(ctx)
			healthResponse, err := drp.Client.Do(req)
			if err != nil {
				return err
			}
			defer healthResponse.Body.Close()

			if healthResponse.StatusCode != http.StatusOK {
				return &retryPath{err: fmt.Errorf("unexpected status: %s", healthResponse.Status)}
			}

			return nil
		})

		// determine whether to retry with another endpoint
		switch t := err.(type) {
		case *retryPath:
			// return the nested error if this is the last ping attempt
			err = t.err
			continue pathLoop
		case kerrors.Aggregate:
			// if any aggregated error indicates a possible retry, do it
			for _, err := range t.Errors() {
				if _, ok := err.(*retryPath); ok {
					continue pathLoop
				}
			}
		}

		break
	}

	return registryURL, err
}

// DryRunRegistryPinger implements RegistryPinger.
type DryRunRegistryPinger struct {
}

// Ping implements Ping method.
func (*DryRunRegistryPinger) Ping(registry string) (*url.URL, error) {
	return url.Parse("https://" + registry)
}

// TryProtocolsWithRegistryURL runs given action with different protocols until no error is returned. The
// https protocol is the first attempt. If it fails and allowInsecure is true, http will be the next. Obtained
// errors will be concatenated and returned.
func TryProtocolsWithRegistryURL(ctx context.Context, registry string, allowInsecure bool, action func(registryURL url.URL) error) (*url.URL, error) {
	errs := []error{}

	if !strings.Contains(registry, "://") {
		registry = "unset://" + registry
	}
	url, err := url.Parse(registry)
	if err != nil {
		return nil, err
	}
	var protos []string
	switch {
	case len(url.Scheme) > 0 && url.Scheme != "unset":
		protos = []string{url.Scheme}
	case allowInsecure || networkutils.IsPrivateAddress(registry):
		protos = []string{HTTPS, HTTP}
	default:
		protos = []string{HTTPS}
	}
	registry = url.Host

	log := logging.FromContext(ctx).Named("TryProtocols")
	for _, proto := range protos {
		log.Debugw("Trying", "protocol", proto, "registry", registry)
		url.Scheme = proto
		err := action(*url)
		if err == nil {
			return url, nil
		}

		log.Debugw("Trying failed", "protocol", proto, "registry", registry, "error", err)
		if _, ok := err.(*errcode.Errors); ok {
			// we got a response back from the registry, so return it
			return url, err
		}

		caErr := errors.Unwrap(err)
		_, isCertificateInvalidError := caErr.(x509.CertificateInvalidError)
		_, isHostnameError := caErr.(x509.HostnameError)

		if (isCertificateInvalidError || isHostnameError) && allowInsecure {
			log.Debugw("Ignored invalid certificate", "protocol", proto, "registry", registry, "error", err)
			return url, nil
		}

		errs = append(errs, err)
		if proto == HTTPS && strings.Contains(err.Error(), "server gave HTTP response to HTTPS client") && !allowInsecure {
			errs = append(errs, fmt.Errorf("\n* Append --force-insecure if you really want to prune the registry using insecure connection"))
		} else if proto == HTTP && strings.Contains(err.Error(), "malformed HTTP response") {
			errs = append(errs, fmt.Errorf("\n* Are you trying to connect to a TLS-enabled registry without TLS?"))
		}
	}

	return nil, kerrors.NewAggregate(errs)
}

// retryPath is an error indicating that another connection attempt may be retried with a different path
type retryPath struct{ err error }

func (rp *retryPath) Error() string { return rp.err.Error() }
