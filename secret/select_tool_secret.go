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
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kclient "github.com/katanomi/pkg/client"
	kerrors "github.com/katanomi/pkg/errors"
	"github.com/katanomi/pkg/namespace"
)

var (
	noSecretSelected = fmt.Errorf("no matching secret based on URL, please specify.")
)

// SelectToolSecretByRefOrLabelOrURL will select secret by secretRef first, if not found, will select by label, if not found, will select by url
func SelectToolSecretByRefOrLabelOrURL(ctx context.Context, currentNamespace string, url string, secretRef *corev1.ObjectReference) (*corev1.Secret, error) {
	clt := kclient.Client(ctx)
	if clt == nil {
		return nil, kerrors.ErrNilPointer
	}
	log := logging.FromContext(ctx)
	ctx = namespace.WithNamespace(ctx, currentNamespace)

	// If secret is specified, try to get directly.
	// Otherwise, select secret by url.
	if secretRef != nil && secretRef.Name != "" {
		secret, err := GetSecretByRefOrLabel(ctx, clt, secretRef)
		if err != nil {
			log.Infow("failed to get secret by ref", "secretRef", secretRef.String(), "err", err)
		}
		return secret, err
	}

	if url == "" {
		return nil, noSecretSelected
	}

	// if it is plugin secret, we should only select secret that synced from integraton directly
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      metav1alpha1.SecretSyncMutationLabelKey,
			Operator: metav1.LabelSelectorOpDoesNotExist,
		}},
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}

	selectOption := SelectSecretOption{
		Namespace:       currentNamespace,
		LabelSelector:   selector,
		PerferredSecret: metav1alpha1.GetNamespacedNameFromRef(secretRef),
		// https://github.com/katanomi/spec/blob/main/docs/core/3.core.credential.selection.md
		Scene: string(metav1alpha1.ResourcePathSceneHttpClone),
	}

	secret, err := SelectToolSecret(log, clt, url, selectOption)
	if err != nil {
		log.Errorw("error to select tool secret", "url", url, "err", err)
		return nil, err
	}
	if secret == nil {
		log.Infow("not select any tool secret", "url", url, "secretRef", secretRef.String())
		return nil, noSecretSelected
	}

	return secret, nil
}
