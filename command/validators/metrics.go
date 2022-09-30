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
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// NewMetric construct a new metric
func NewMetric(metrics map[string]string) *Metric {
	return &Metric{
		metrics: metrics,
	}
}

// Metric help to validate metrics
type Metric struct {
	metrics map[string]string
}

// GetMetricValue get metric value
func (p *Metric) GetMetricValue(metric string) (value string, exist bool) {
	if p.metrics == nil {
		return "", false
	}
	value, exist = p.metrics[metric]
	return
}

// ValidateInt help to validate int metrics
func (p *Metric) ValidateInt(base *field.Path, metric string, min, max *int) (errs field.ErrorList) {
	value, exist := p.GetMetricValue(metric)
	if !exist {
		return nil
	}

	path := base.Child(metric)
	v, rateError := strconv.Atoi(strings.TrimSpace(value))
	if rateError != nil {
		errs = append(errs, field.Invalid(path, value, `value is not a number and cannot be parsed`))
		return errs
	}

	if min != nil && v < *min {
		errs = append(errs, field.Invalid(path, value, fmt.Sprintf("value should be greater than %d", *min)))
	}

	if max != nil && v > *max {
		errs = append(errs, field.Invalid(path, value, fmt.Sprintf("value should be less than %d", *max)))
	}

	return
}

// ValidateFloat help to validate float metrics
func (p *Metric) ValidateFloat(base *field.Path, metric string, min, max *float64) (errs field.ErrorList) {
	value, exist := p.GetMetricValue(metric)
	if !exist {
		return nil
	}

	path := base.Child(metric)
	v, rateError := strconv.ParseFloat(strings.TrimRight(value, "%"), 64)
	if rateError != nil {
		errs = append(errs, field.Invalid(base.Child(metric), value, `value is not a float number and cannot be parsed`))
		return errs
	}
	if min != nil && v < *min {
		errs = append(errs, field.Invalid(path, value, fmt.Sprintf("value should be greater than %g", *min)))
	}
	if max != nil && v > *max {
		errs = append(errs, field.Invalid(path, value, fmt.Sprintf("value should be less than %g", *max)))
	}
	return
}
