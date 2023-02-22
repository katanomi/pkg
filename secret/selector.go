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

// Package secret contains useful functionality for select secret
package secret

import (
	"context"
	"fmt"
	neturl "net/url"
	"sort"
	"strings"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// A better secret selection solution that do not dependency on tool resource about tool secret
// https://github.com/katanomi/spec/blob/main/3.core.credential.selection.md

// SelectSecretOption encapsulate the configuration related to selecting secrets
type SelectSecretOption struct {

	// Scene indicates resource url format in different scenario
	Scene string

	// PerferredSecret will return the secret if it is be selected
	PerferredSecret types.NamespacedName

	// ExcludedSecretTypes exclude some secret types when do selecting
	ExcludedSecretTypes SecretTypeList

	// SecretTypes means only secret which type exist in this list will be selected
	// if it is empty means there is no limit for secret type when selecting
	SecretTypes SecretTypeList

	// Namespace indicates current namespace that current resource belongs.
	// the secret will be searched in this namespace
	// as a default action, secret in the same namespace could be used by other resources in same namespace
	Namespace string
	// GlobalCredentialsNamespace is the namespace that we save global credentials.
	// it is not public to all namespace only until it is bind to one project(namespace)
	// the secret will be searched in this namespace
	GlobalCredentialsNamespace string

	// LabelSelector is label selector when select secret, default will be everything
	LabelSelector labels.Selector

	// IncludeAnnotaion if specified, it needs to be filtered base IncludeAnnotaion. Currently only key value filtering is used.
	IncludeAnnotaion map[string]string
}

type SecretTypeList []corev1.SecretType

func (s SecretTypeList) Contains(e corev1.SecretType) bool {
	for _, secret := range s {
		if secret == e {
			return true
		}
	}
	return false
}

// NewSecretSelectOption just construct SecretSelectOption
func NewSecretSelectOption(preferredSecret types.NamespacedName, namespace string, globalCredentialsNamespace string) (option SelectSecretOption) {

	return SelectSecretOption{
		PerferredSecret:            preferredSecret,
		Namespace:                  namespace,
		GlobalCredentialsNamespace: globalCredentialsNamespace}
}

// SelectToolSecret will select secret according to tool address and resource scope on secret
//
//	clientI could be sigs.k8s.io/controller-runtime/pkg/client.Client or k8s.io/client-go/kubernetes.Interface or "k8s.io/client-go/informers".SharedInformerFactory
//	resourceURL refers resource url lik git http url or harbor http url. eg. https://github.com/example or  build.harbor.com/example
//	namespaces refers to all namespaces where the secret may exist
//
// if no secret was found, secret will be nil and err is nil
// if any errors occurred, err will not be nil
// The meaning of preferred is that if there is a secret with this name, it will be selected first, and preferredNs means the namespace where the secret is located
func SelectToolSecret(logger *zap.SugaredLogger, clientI interface{}, resourceURL string, option SelectSecretOption) (secret *corev1.Secret, err error) {

	logger = logger.Named("secret-selector").With("resourceURL", resourceURL, "option", option)

	defer func() {
		if err != nil {
			logger.Errorw("selected tool secret error", "err", err)
		} else {
			if secret == nil {
				logger.Debugw("selected tool secret is nil")
			} else {
				logger.Debugw("selected tool secret", "secret", secret.Namespace+"/"+secret.Name)
			}
		}
	}()

	var listSecretFunc = func(ns string) (*corev1.SecretList, error) {

		var secretList = &corev1.SecretList{}

		switch client := clientI.(type) {
		case ctrlclient.Client:
			var listOpts = []ctrlclient.ListOption{}
			if ns != "" {
				listOpts = append(listOpts, ctrlclient.InNamespace(ns))
			}
			if option.LabelSelector != nil {
				listOpts = append(listOpts, ctrlclient.MatchingLabelsSelector{Selector: option.LabelSelector})
			}

			err := client.List(context.Background(), secretList, listOpts...)
			return secretList, err
		case kubernetes.Interface:
			listOpts := metav1.ListOptions{ResourceVersion: "0"}
			if option.LabelSelector != nil {
				listOpts.LabelSelector = option.LabelSelector.String()
			}
			secretList, err := client.CoreV1().Secrets(ns).List(context.Background(), listOpts)
			return secretList, err
		case k8sinformers.SharedInformerFactory:
			selector := labels.Everything()
			if option.LabelSelector != nil {
				selector = option.LabelSelector
			}
			list, err := client.Core().V1().Secrets().Lister().Secrets(ns).List(selector)
			if err != nil {
				return secretList, err
			}
			secretList.Items = make([]corev1.Secret, 0, len(list))
			for _, item := range list {
				secretList.Items = append(secretList.Items, *item)
			}
			return secretList, nil
		default:
			return secretList, fmt.Errorf("error type of client : %T", clientI)
		}
	}

	var secretList, globalSecretList *corev1.SecretList = &corev1.SecretList{}, &corev1.SecretList{}

	secretList, err = listSecretFunc(option.Namespace)
	if err != nil {
		return nil, err
	}
	if option.GlobalCredentialsNamespace != "" {
		globalSecretList, err = listSecretFunc(option.GlobalCredentialsNamespace)
		if err != nil {
			return nil, err
		}
	}

	secret, err = selectToolSecret(logger, secretList.Items, globalSecretList.Items, resourceURL, option)
	return
}

func sortSecretList(secrets []corev1.Secret) []corev1.Secret {
	var sortedList = SortedSecretList(secrets)
	sort.Sort(sortedList)
	return sortedList
}

func selectToolSecret(logger *zap.SugaredLogger, secretList []corev1.Secret, globalSecretList []corev1.Secret, resourceURL string, option SelectSecretOption) (secret *corev1.Secret, err error) {

	correctResUrl := resourceURL
	if !strings.HasPrefix(resourceURL, "http://") && !strings.HasPrefix(resourceURL, "https://") {
		// ensure resource url is valid neturl
		correctResUrl = "http://" + correctResUrl
	}

	resourceU, err := neturl.Parse(correctResUrl)
	if err != nil {
		return nil, err
	}

	usableSecrets := []corev1.Secret{}
	usableSecrets = append(usableSecrets, selectToolSecretFrom(logger, secretList, false, resourceU, option)...)
	usableSecrets = append(usableSecrets, selectToolSecretFrom(logger, globalSecretList, true, resourceU, option)...)
	if len(usableSecrets) == 0 {
		return nil, nil
	}
	find, secretIndex := findPreferredSecret(usableSecrets, option.PerferredSecret.Namespace, option.PerferredSecret.Name)
	if find {
		return &usableSecrets[secretIndex], nil
	}

	// get the latest secret.
	if newestIndex := findNewestSecret(usableSecrets); newestIndex > -1 {
		return &usableSecrets[newestIndex], nil
	}

	return &usableSecrets[0], nil
}

func selectToolSecretFrom(logger *zap.SugaredLogger, secretList []corev1.Secret, isGlobal bool, resourceURL *neturl.URL, option SelectSecretOption) []corev1.Secret {
	usableSecrets := make([]corev1.Secret, 0)

	for _, _secret := range secretList {
		var sec = _secret
		logger := logger.With("secret", sec.Namespace+"/"+sec.Name)

		if len(option.SecretTypes) != 0 && !option.SecretTypes.Contains(sec.Type) {
			logger.Debugw("secret type mismatch")
			continue
		}
		if len(option.ExcludedSecretTypes) != 0 && option.ExcludedSecretTypes.Contains(sec.Type) {
			logger.Debugw("secret type is excluded")
			continue
		}

		address := sec.Annotations[metav1alpha1.IntegrationAddressAnnotation]
		if address == "" {
			logger.Debugw("secret address annotation is empty")
			continue
		}

		if !HasAnnotationsKey(sec.Annotations, option.IncludeAnnotaion, false) {
			logger.Debugw("secret not macth expected annotation")
			continue
		}

		toolAddress, err := neturl.Parse(address)
		if err != nil {
			logger.Infow("integration address is invalid", "address", address)
			continue
		}

		if resourceURL.Host != toolAddress.Host {
			logger.Debugw("skip select secret, host not same", "resourceHost", resourceURL.Host, "annotationAddressHost", toolAddress.Host)
			continue
		}

		scopes := sec.Annotations[metav1alpha1.IntegrationResourceScope]
		if scopes == "" {
			logger.Debugw("skip select secret, scopes is empty")
			continue
		}

		if isGlobal {
			if sec.Annotations[metav1alpha1.IntegrationSecretApplyNamespaces] == "" {
				logger.Debugw("secret apply namespaces is empty")
				continue
			}
			applyNamespaces := strings.Split(sec.Annotations[metav1alpha1.IntegrationSecretApplyNamespaces], ",")
			if !containsInArray(applyNamespaces, option.Namespace) {
				logger.Debugw("skip select global secret, apply namespace is not contains target namespace", "applyNamespaces", applyNamespaces, "namespace", option.Namespace)
				continue
			}
		}

		logger.Debugw("secret scopes", "scopes", scopes)
		scopeItems := strings.Split(scopes, ",")
		pathFmtJson := sec.Annotations[metav1alpha1.IntegrationSecretResourcePathFmt]
		subPathFmtJson := sec.Annotations[metav1alpha1.IntegrationSecretSubResourcePathFmt]
		resourcePathFormat := NewResourcePathFormat(pathFmtJson, subPathFmtJson)
		for _, scope := range scopeItems {
			if scope == "" {
				continue
			}

			acceptableResourceUrl := resourcePathFormat.FormatPathByScene(metav1alpha1.ResourcePathScene(option.Scene), scope)

			inputResourceURL := strings.ToLower(resourceURL.Path)
			if strings.HasSuffix(inputResourceURL, ".git") {
				inputResourceURL = strings.TrimSuffix(inputResourceURL, ".git")
			}
			if !strings.HasSuffix(inputResourceURL, "/") {
				inputResourceURL = inputResourceURL + "/"
			}
			if acceptableResourceUrl == "/" || strings.HasPrefix(inputResourceURL, strings.ToLower(acceptableResourceUrl)) {
				usableSecrets = append(usableSecrets, sec)
				break
			}
		}
	}

	return usableSecrets
}

func findPreferredSecret(secrets []corev1.Secret, preferredNs, preferred string) (find bool, index int) {
	if preferred != "" && preferredNs != "" {
		for index, secret := range secrets {
			if secret.Name == preferred && preferredNs == secret.Namespace {
				return true, index
			}
		}
	}
	return false, -1
}

func findNewestSecret(secrets []corev1.Secret) (index int) {
	newestTime := metav1.Time{}
	newestIndex := -1
	for index := range secrets {
		if !secrets[index].CreationTimestamp.Before(&newestTime) {
			newestTime = secrets[index].CreationTimestamp
			newestIndex = index
		}
	}

	return newestIndex
}

func containsInArray(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

// HasAnnotationsKey whether to include the expected annotation key, returns true when the expected annotation is empty.
func HasAnnotationsKey(annotations, expectedAnnotaions map[string]string, matchValue bool) bool {
	if expectedAnnotaions == nil {
		return true
	}

	for k, expectedValue := range expectedAnnotaions {
		if annotations == nil {
			return false
		}

		annotationValue, ok := annotations[k]
		if !ok {
			return false
		}

		if matchValue && expectedValue != annotationValue {
			return false
		}
	}

	return true
}
