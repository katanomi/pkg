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

package sharedmain

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/google/go-cmp/cmp"
	"github.com/katanomi/pkg/apis/meta/v1alpha1"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
)

type MockPlugin struct {
	VersionAttributes map[string]map[string][]string
	DataAttributes    map[string][]string
	DisplayColumnsMap map[string]v1alpha1.DisplayColumns
}

func (m *MockPlugin) GetVersionAttributes(version string) map[string][]string {
	return m.VersionAttributes[version]
}

func (m *MockPlugin) SetVersionAttributes(version string, attributes map[string][]string) {
	m.VersionAttributes[version] = attributes
}

func (m *MockPlugin) SetAttribute(k string, values ...string) {
	m.DataAttributes[k] = values
}

func (m *MockPlugin) GetAttribute(k string) []string {
	return m.DataAttributes[k]
}

func (m *MockPlugin) Attributes() map[string][]string {
	return m.DataAttributes
}

func (m *MockPlugin) GetDisplayColumns() map[string]v1alpha1.DisplayColumns {
	return m.DisplayColumnsMap
}

func (m *MockPlugin) SetDisplayColumns(k string, values ...v1alpha1.DisplayColumn) {
	m.DisplayColumnsMap[k] = values
}

func NewMockPlugin() *MockPlugin {
	return &MockPlugin{
		VersionAttributes: map[string]map[string][]string{},
		DataAttributes:    map[string][]string{},
		DisplayColumnsMap: map[string]v1alpha1.DisplayColumns{},
	}
}

func TestAppBuilder_PluginVersionAttributes(t *testing.T) {
	t.Run("read version attributes to plugin", func(t *testing.T) {
		g := NewGomegaWithT(t)
		logger, _ := zap.NewDevelopment()
		app := AppBuilder{
			Logger: logger.Sugar(),
		}
		plugin := NewMockPlugin()

		app.PluginVersionAttributes(plugin, "testdata/versionattributes.yaml")

		diff := cmp.Diff(plugin.VersionAttributes["version"], map[string][]string{"attr": {"field1", "field2"}})
		g.Expect(diff).To(BeEmpty())
	})
}

func TestAppBuilder_PluginAttributes(t *testing.T) {
	g := NewGomegaWithT(t)
	logger, _ := zap.NewDevelopment()
	app := AppBuilder{
		Logger: logger.Sugar(),
	}
	plugin := NewMockPlugin()

	app.PluginAttributes(plugin, "testdata/attributes.yaml")

	diff := cmp.Diff(plugin.DataAttributes["attr"], []string{"field1", "field2"})
	g.Expect(diff).To(BeEmpty())
}

func TestAppBuilder_PluginDisplayColumns(t *testing.T) {
	t.Run("read display columns to plugin", func(t *testing.T) {
		g := NewGomegaWithT(t)
		logger, _ := zap.NewDevelopment()
		app := AppBuilder{
			Logger: logger.Sugar(),
		}
		plugin := NewMockPlugin()

		app.PluginDisplayColumns(plugin, "testdata/displaycolumns.yaml")

		diff := cmp.Diff(plugin.DisplayColumnsMap["projectColumns"], v1alpha1.DisplayColumns{{
			Name:        "name",
			Field:       "metadata.name",
			DisplayName: "_.integrations.project.columns.name",
			Properties: &runtime.RawExtension{
				Raw: []byte(`{"width":"200"}`),
			},
		}})
		g.Expect(diff).To(BeEmpty())
	})
}

func TestAppWithFieldIndexer(t *testing.T) {
	t.Run("append one field indexer ", func(t *testing.T) {
		g := NewGomegaWithT(t)
		a := App("test").WithFieldIndexer(FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(1))
	})
	t.Run("append more than one field indexer", func(t *testing.T) {
		g := NewGomegaWithT(t)

		a := App("test").WithFieldIndexer(FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.key",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["key"]}
			},
		}).WithFieldIndexer(FieldIndexer{
			Obj:   &corev1.ConfigMap{},
			Field: "data.name",
			ExtractValue: func(object client.Object) []string {
				return []string{object.(*corev1.ConfigMap).Data["name"]}
			},
		})
		g.Expect(a.fieldIndexeres).Should(HaveLen(2))
	})
}
