/*
Copyright 2024 The AlaudaDevops Authors.

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

package testing

import (
	"encoding/base64"

	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// AddReportEntriesAsYaml adds multiple report entries formatted as YAML
// to the current SpecReport.
// The Report entry visibility is set to ReportEntryVisibilityFailureOrVerbose.
//
// Example:
//
//	configMap := &corev1.ConfigMap{}
//	service := &corev1.Service{}
//	client.Get(ctx, key, configMap)
//	client.Get(ctx, serviceKey, service)
//	AddReportEntriesAsYaml("configuration", configMap, service)
//
// Note: Any pointer will be only printed as the final value when testing is completed
// for current value provide the object as value:
//
//	AddReportEntriesAsYaml("previous-configuration", *configMap)
func AddReportEntriesAsYaml(name string, objs ...interface{}) {
	for _, obj := range objs {
		AddReportEntryAsYaml(name, obj)
	}
}

// AddReportEntryAsYaml marshals the given object as YAML and adds it
// to the test report with the given name.
// The Report entry visibility is set to ReportEntryVisibilityFailureOrVerbose.
// Note: Any pointer will be only printed as the final value when testing is completed
// for current value provide the object as value.
//
// Example:
//
//	configMap := &corev1.ConfigMap{}
//	client.Get(ctx, key, configMap)
//	AddReportEntryAsYaml("configuration", configMap)
//
// Note: Any pointer will be only printed as the final value when testing is completed
// for current value provide the object as value:
//
//	AddReportEntryAsYaml("previous-configuration", *configMap)
//
// Will automatically redact secret data for corev1.Secret objects.
func AddReportEntryAsYaml(name string, obj interface{}) {
	obj = redactSecrets(obj)
	if bts, err := yaml.Marshal(obj); err == nil {
		AddReportEntry(name, ReportEntryVisibilityFailureOrVerbose, string(bts))
	} else {
		AddReportEntry(name, ReportEntryVisibilityFailureOrVerbose, obj)
	}
}

const redactedString = "REDACTED"

var redactedStringBase64 = base64.StdEncoding.EncodeToString([]byte(redactedString))

// redactSecrets redacts sensitive fields in Kubernetes Secrets and SecretLists.
// It returns a deep copy of the input with secret data replaced with
// redacted values.
func redactSecrets(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	switch obj := obj.(type) {
	case corev1.Secret:
		objCopy := obj.DeepCopy()
		objCopy.Data = redactSecretData(obj.Data)
		return *objCopy
	case *corev1.Secret:
		objCopy := obj.DeepCopy()
		objCopy.Data = redactSecretData(obj.Data)
		return objCopy
	case *corev1.SecretList:
		objCopy := obj.DeepCopy()
		for i := range objCopy.Items {
			objCopy.Items[i] = redactSecrets(objCopy.Items[i]).(corev1.Secret)
		}
		return objCopy
	case corev1.SecretList:
		objCopy := obj.DeepCopy()
		for i := range objCopy.Items {
			objCopy.Items[i] = redactSecrets(objCopy.Items[i]).(corev1.Secret)
		}
		return *objCopy
	default:
		// fallthrough returns the original object
		return obj
	}
}

// redactSecretData redacts the sensitive data in a Secret's data map
// by replacing the values with a redacted string. It returns a new
// map with the redacted values.
func redactSecretData(data map[string][]byte) map[string][]byte {
	redactedData := map[string][]byte{}
	for k := range data {
		redactedData[k] = []byte(redactedStringBase64)
	}
	return redactedData
}
