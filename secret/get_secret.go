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

package secret

import (
	"context"
	"fmt"

	pkgnamespace "github.com/katanomi/pkg/namespace"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	kerrors "github.com/katanomi/pkg/errors"
)

// GetSecretByRefOrLabel retrieves an secret
// If the ref.namespace is empty, it will use the namespace from the ctx.
// If the resource does not exist, it will try to match according to the label.
func GetSecretByRefOrLabel(ctx context.Context, clt client.Client, ref *corev1.ObjectReference) (obj *corev1.Secret, err error) {
	if clt == nil || ref == nil {
		return nil, kerrors.ErrNilPoint
	}

	log := logging.FromContext(ctx)
	obj = &corev1.Secret{}
	objKey := client.ObjectKey{Namespace: ref.Namespace, Name: ref.Name}
	if objKey.Namespace == "" {
		objKey.Namespace = pkgnamespace.NamespaceValue(ctx)
	}
	if err = clt.Get(ctx, objKey, obj); err == nil {
		return obj, nil
	} else if !apierrors.IsNotFound(err) {
		return nil, err
	}

	// not found secret, we should find secret by labels
	matchingLabels := map[string]string{
		metav1alpha1.NamespaceLabelKey: objKey.Namespace,
		metav1alpha1.SecretLabelKey:    objKey.Name,
	}

	// label selector to select secret matching ns and name matched and the original secret
	labelSelector := &metav1.LabelSelector{
		MatchLabels: matchingLabels,
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      metav1alpha1.SecretSyncMutationLabelKey,
			Operator: metav1.LabelSelectorOpDoesNotExist,
		}},
	}
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}

	secretList := &corev1.SecretList{}
	if err = clt.List(ctx, secretList, client.InNamespace(objKey.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return nil, err
	}

	if len(secretList.Items) == 0 {
		log.Errorw("not found any secret by labels", "labels", matchingLabels, "secretRef", ref)
		err = fmt.Errorf("secret %s not exist or matching by labels '%v'", objKey.String(), matchingLabels)
		return nil, err
	}
	if len(secretList.Items) > 1 {
		log.Infow("found multiple secrets by labels", "labels", matchingLabels, "#secretList.Items", len(secretList.Items))
	}
	return &secretList.Items[0], nil
}
