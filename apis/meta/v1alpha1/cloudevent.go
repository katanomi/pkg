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

package v1alpha1

import (
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudEvent struct {
	ID      string `json:"id,omitempty"`
	Source  string `json:"source,omitempty"`
	Subject string `json:"subject,omitempty"`
	// Type of event
	Type string `json:"type,omitempty"`
	// Data event payload
	Data            string            `json:"data,omitempty"`
	Time            metav1.Time       `json:"time,omitempty"`
	SpecVersion     string            `json:"specversion,omitempty"`
	DataContentType string            `json:"datacontenttype,omitempty"`
	Extensions      map[string]string `json:"extensions,omitempty"`
}

func (evt *CloudEvent) From(event cloudevents.Event) *CloudEvent {
	evt.ID = event.ID()
	evt.Source = event.Source()
	evt.Data = string(event.Data())
	evt.Subject = event.Subject()
	evt.DataContentType = event.DataContentType()
	evt.Type = event.Type()
	evt.SpecVersion = event.SpecVersion()
	evt.Time = metav1.NewTime(event.Time())
	for key, val := range event.Extensions() {
		if evt.Extensions == nil {
			evt.Extensions = map[string]string{}
		}

		var str string

		switch v := val.(type) {
		case string:
			str = v
		case int, int8, int16, int32, int64:
			str = fmt.Sprintf("%d", v)
		case time.Time:
			str = time.Time(v).String()
		case float32, float64:
			str = fmt.Sprintf("%d", v)
		case bool:
			str = fmt.Sprintf("%t", v)
		default:
			bts, _ := json.Marshal(v)
			str = string(bts)
		}
		evt.Extensions[key] = str
	}
	return evt
}
