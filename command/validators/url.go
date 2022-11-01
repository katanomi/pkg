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
	"net/url"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateURLFunc function to validate url
type ValidateURLFunc func(url *url.URL) (ok bool, errMsg string)

// NewURL construct a new metric
func NewURL() *URL {
	return &URL{}
}

// URL help to validate url
type URL struct {
	errMsg   string
	validate ValidateURLFunc
}

// SetErrMsg customize error message
func (p *URL) SetErrMsg(msg string) *URL {
	p.errMsg = msg
	return p
}

// SetValidate customize validate function
func (p *URL) SetValidate(f ValidateURLFunc) *URL {
	p.validate = f
	return p
}

// Validate Support to verify multiple urls at the same time
func (p *URL) Validate(path *field.Path, urls ...string) (errs field.ErrorList) {
	for idx, _url := range urls {
		u, err := url.Parse(_url)
		if err != nil {
			errs = append(errs, field.Invalid(path.Index(idx), _url, p.getErrMsg()))
			continue
		}
		if p.validate == nil {
			continue
		}
		if ok, errMsg := p.validate(u); !ok {
			errs = append(errs, field.Invalid(path.Index(idx), _url, p.getErrMsg(errMsg)))
		}
	}
	return errs
}

// getErrMsg get error message, if not set, use default error message
func (p *URL) getErrMsg(msgs ...string) string {
	for _, msg := range msgs {
		if msg != "" {
			return msg
		}
	}
	if p.errMsg != "" {
		return p.errMsg
	}
	return "invalid url"
}
