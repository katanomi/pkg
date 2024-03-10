/*
Copyright 2024 The Katanomi Authors.

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
	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/format"
	uberzap "go.uber.org/zap"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// InitializeGinkgoConfig initializes the Ginkgo configuration.
func InitializeGinkgoConfig() {

	// Disable the string length limit when Gomega outputs information in tests.
	// The default is 4000, which may result in incomplete log output.
	// Ref: https://onsi.github.io/gomega/#adjusting-output
	format.MaxLength = 0
}

// NewGinkgoLogger creates a new logger for Ginkgo tests.
func NewGinkgoLogger() *uberzap.SugaredLogger {
	return zap.NewRaw(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)).Sugar()
}

// GetDefaultLogger returns the default logger for testing
func GetDefaultLogger() *uberzap.SugaredLogger {
	return NewGinkgoLogger()
}

// InitGinkgoWithLogger initializes Ginkgo and returns a logger.
func InitGinkgoWithLogger() *uberzap.SugaredLogger {
	InitializeGinkgoConfig()

	// set the logger for the controller-runtime package.
	logger := NewGinkgoLogger()
	logf.SetLogger(zapr.NewLogger(logger.Desugar()))

	return NewGinkgoLogger()
}
