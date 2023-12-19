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

package v2

import (
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var testClient = resty.New()
var mate = base.Meta{Version: "v1.3.4", BaseURL: "https://plugin.com"}
var secret = corev1.Secret{
	Type: corev1.SecretTypeBasicAuth,
	Data: map[string][]byte{"username": []byte("username")},
}
var pluginUrl, _ = apis.ParseURL("https://example.com/")
var pluginClient = NewPluginClientV2(&duckv1.Addressable{URL: pluginUrl}, mate, secret, base.ClientOpts(testClient))
var listOption = metav1alpha1.ListOptions{
	Page:         2,
	ItemsPerPage: 3,
}

var _ = BeforeSuite(func() {
	httpmock.ActivateNonDefault(testClient.GetClient())
})

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

func TestPluginClientV2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Client V2 Suite")
}

func diff(a, b interface{}) string {
	return cmp.Diff(a, b)
}

var ignoredFakeFieldTypes = []reflect.Type{
	reflect.TypeOf(duckv1.Addressable{}),
	reflect.TypeOf(metav1.Time{}),
	reflect.TypeOf(runtime.RawExtension{}),
}

func isIgnoredField(t reflect.Type) bool {
	for _, ignoredType := range ignoredFakeFieldTypes {
		if ignoredType.Name() == t.Name() && ignoredType.PkgPath() == t.PkgPath() {
			return true
		}
	}
	return false
}

// fakeValue generate random test data for string and int field only,
// that's sufficient for testing
func fakeValue(v reflect.Value) {
	t := v.Type()
	switch t.Kind() {
	case reflect.Ptr:
		if isIgnoredField(t.Elem()) {
			return
		}
		if v.IsNil() {
			nv := reflect.New(t.Elem())
			fakeValue(nv.Elem())
			v.Set(nv)
		} else {
			fakeValue(v.Elem())
		}
	case reflect.String:
		v.SetString(rand.String(10))
	case reflect.Int:
		v.SetInt(int64(rand.Intn(100)))
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if !v.Field(i).CanSet() || isIgnoredField(v.Field(i).Type()) {
				continue
			}
			fakeValue(v.Field(i))
		}
	default:
		// nothing
	}
}

// fakeStruct fills a struct with random values
func fakeStruct[T interface{}]() *T {
	t := new(T)
	fakeValue(reflect.ValueOf(t).Elem())
	return t
}
