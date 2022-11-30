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

package validators

import (
	"testing"

	"github.com/katanomi/pkg/pointer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ = Describe("Test.Metrics.GetMetricValue", func() {
	var metric *Metric
	Context("metrics is nil", func() {
		BeforeEach(func() {
			metric = NewMetric(nil)
		})

		When("get metric value", func() {
			It("should return empty string and false", func() {
				value, exist := metric.GetMetricValue("test")
				Expect(value).To(Equal(""))
				Expect(exist).To(BeFalse())
			})
		})
	})

	Context("metrics is not nil", func() {
		BeforeEach(func() {
			metric = NewMetric(map[string]string{
				"test": "value",
			})
		})

		When("get an exist metric", func() {
			It("should return the value and true", func() {
				value, exist := metric.GetMetricValue("test")
				Expect(value).To(Equal("value"))
				Expect(exist).To(BeTrue())
			})
		})

		When("get a not exist metric", func() {
			It("should return empty string and false", func() {
				value, exist := metric.GetMetricValue("not-exist")
				Expect(value).To(Equal(""))
				Expect(exist).To(BeFalse())
			})
		})
	})
})

func TestMetric_ValidateFloat(t *testing.T) {
	g := NewGomegaWithT(t)
	path := field.NewPath("root")
	tests := []struct {
		metrics  map[string]string
		metric   string
		min      *float64
		max      *float64
		wantErrs bool
	}{
		{
			metrics:  nil,
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			metric:   "not-exist",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1.1",
			},
			metric:   "test",
			wantErrs: true,
		},
		{
			metrics: map[string]string{
				"test": "10.1%%",
			},
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			min:      pointer.Float64(1),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			min:      pointer.Float64(10.1),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			min:      pointer.Float64(11),
			metric:   "test",
			wantErrs: true,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			max:      pointer.Float64(11),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			max:      pointer.Float64(10.1),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10.1",
			},
			max:      pointer.Float64(0),
			metric:   "test",
			wantErrs: true,
		},
	}
	for _, tt := range tests {
		p := &Metric{metrics: tt.metrics}
		errs := p.ValidateFloat(path, tt.metric, tt.min, tt.max)
		g.Expect(len(errs) > 0).To(Equal(tt.wantErrs))
	}
}

func TestMetric_ValidateInt(t *testing.T) {
	g := NewGomegaWithT(t)
	path := field.NewPath("root")
	tests := []struct {
		metrics  map[string]string
		metric   string
		min      *int
		max      *int
		wantErrs bool
	}{
		{
			metrics:  nil,
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10x",
			},
			metric:   "test",
			wantErrs: true,
		},
		{
			metrics: map[string]string{
				"test": "10  ",
			},
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			metric:   "not-exist",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			min:      pointer.Int(1),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			min:      pointer.Int(10),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			min:      pointer.Int(11),
			metric:   "test",
			wantErrs: true,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			max:      pointer.Int(11),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			max:      pointer.Int(10),
			metric:   "test",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"test": "10",
			},
			max:      pointer.Int(0),
			metric:   "test",
			wantErrs: true,
		},
	}
	for _, tt := range tests {
		p := &Metric{metrics: tt.metrics}
		errs := p.ValidateInt(path, tt.metric, tt.min, tt.max)
		g.Expect(len(errs) > 0).To(Equal(tt.wantErrs))
	}
}

func TestMetric_ValidateStringEnums(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		metrics  map[string]string
		metric   string
		enums    []string
		wantErrs bool
	}{
		{
			metrics: map[string]string{
				"key": "value",
			},
			metric:   "not-exist-key",
			wantErrs: false,
		},
		{
			metrics: map[string]string{
				"key": "value",
			},
			metric:   "key",
			enums:    []string{"val1", "val2"},
			wantErrs: true,
		},
		{
			metrics: map[string]string{
				"key": "value",
			},
			metric:   "key",
			enums:    []string{"val1", "val2", "value"},
			wantErrs: false,
		},
	}
	testPath := field.NewPath("test")
	for _, tt := range tests {
		p := &Metric{
			metrics: tt.metrics,
		}
		errs := p.ValidateStringEnums(testPath, tt.metric, tt.enums...)
		g.Expect(len(errs) > 0).To(Equal(tt.wantErrs))
	}
}
