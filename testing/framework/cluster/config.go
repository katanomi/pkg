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

package cluster

import (
	"context"
	"strings"
	"sync"

	"github.com/katanomi/pkg/testing"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	e2eConfigNSKey   = "E2E_CONFIG_NAMESPACE"
	e2eConfigNameKey = "E2E_CONFIG_NAME_PREFIX"
)

var (
	e2eConfigNs         = testing.GetDefaultEnv(e2eConfigNSKey, "katanomi-e2e")
	e2eConfigNamePrefix = testing.GetDefaultEnv(e2eConfigNameKey, "e2e-config")
)

// GetConfigFromContext get the config information from context
func GetConfigFromContext(ctx context.Context) interface{} {
	return ctx.Value(configCondition{})
}

// NewConfigCondition construct an configCondition object
// `obj` is a pointer which used to unmarshal configuration to
func NewConfigCondition(configName string, obj interface{}) *configCondition {
	c := &configCondition{
		name: configName,
		obj:  obj,
	}

	return c
}

type configCondition struct {
	name string
	obj  interface{}
}

// Condition implement the Condition interface
func (c *configCondition) Condition(testCtx *TestContext) error {
	configData, err := NewE2EConfig(c.name).GetConfig(testCtx.Context, testCtx.Client)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(configData), c.obj); err != nil {
		return err
	}

	testCtx.Context = context.WithValue(testCtx.Context, configCondition{}, c.obj)
	return nil
}

// NewE2EConfig construct a `TestConfig` with special configmap name
func NewE2EConfig(cmName string) *TestConfig {
	name := []string{strings.TrimSuffix(e2eConfigNamePrefix, "-"), cmName}
	return &TestConfig{
		Namespace: e2eConfigNs,
		Name:      strings.Join(name, "-"),
	}
}

// TestConfig Provide a unified way to get configuration from configmap
type TestConfig struct {
	Namespace string
	Name      string

	lock sync.RWMutex
	once sync.Once

	data string
}

func (c *TestConfig) initData(ctx context.Context, clt client.Client) error {
	var err error
	c.once.Do(func() {
		key := types.NamespacedName{Namespace: c.Namespace, Name: c.Name}
		cm := &v1.ConfigMap{}
		if err = clt.Get(ctx, key, cm); err != nil {
			return
		}
		if cm.Data != nil {
			c.data = cm.Data["config"]
		}
	})
	return err
}

// GetConfig get configuration by specified client
func (c *TestConfig) GetConfig(ctx context.Context, clt client.Client) (string, error) {
	if err := c.initData(ctx, clt); err != nil {
		return "", err
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data, nil
}
