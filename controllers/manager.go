/*
Copyright 2023 The AlaudaDevops Authors.

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

package controllers

import (
	"fmt"

	"k8s.io/utils/strings"

	"k8s.io/apimachinery/pkg/runtime"
	cliengorecord "k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	// MessageShortenLength default message shorten length in event
	MessageShortenLength = 1000
)

// ControllerManager is just like sigs.k8s.io/controller-runtime.Manager
// and decorate GetEventRecorderFor function to shorten message length
type ControllerManager struct {
	ctrl.Manager
}

// GetEventRecorderFor is just like sigs.k8s.io/controller-runtime.Manager.GetEventRecorderFor
// but will return EventRecorder that could shorten message length
func (m ControllerManager) GetEventRecorderFor(name string) cliengorecord.EventRecorder {
	return &EventRecorder{
		EventRecorder:     m.Manager.GetEventRecorderFor(name),
		MessageShortenLen: MessageShortenLength,
	}
}

// EventRecorder is just like k8s.io/client-go/tools/record.EventRecorder, but will shorten message length according to MessageShortenLen
type EventRecorder struct {
	cliengorecord.EventRecorder

	// MessageShortenLen indicates max event message length, it will be shortened if length is more than MessageShortenLen
	MessageShortenLen int
}

// Event is same as k8s.io/client-go/tools/record.EventRecorder.Event, just shorten message filed
func (recorder *EventRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	recorder.EventRecorder.Event(object, eventtype, reason, strings.ShortenString(message, recorder.MessageShortenLen))
}

// Eventf is same as k8s.io/client-go/tools/record.EventRecorder.Eventf, just shorten args filed if the field type is Stringer or Error
func (recorder *EventRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	shortens := recorder.shortenArgs(args...)
	recorder.EventRecorder.Eventf(object, eventtype, reason, messageFmt, shortens...)
}

// AnnotatedEventf is same as k8s.io/client-go/tools/record.EventRecorder.Eventf, just shorten args filed if the field type is Stringer or Error
func (recorder *EventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	shortens := recorder.shortenArgs(args...)
	recorder.EventRecorder.AnnotatedEventf(object, annotations, eventtype, reason, messageFmt, shortens...)
}

func (recorder *EventRecorder) shortenArgs(args ...interface{}) []interface{} {
	shorten := make([]interface{}, 0, len(args))
	for _, item := range args {
		if arg, ok := item.(fmt.Stringer); ok {
			shorten = append(shorten, strings.ShortenString(arg.String(), recorder.MessageShortenLen))
			continue
		}

		if arg, ok := item.(string); ok {
			shorten = append(shorten, strings.ShortenString(arg, recorder.MessageShortenLen))
			continue
		}

		if arg, ok := item.(interface {
			Error() string
		}); ok {
			shorten = append(shorten, strings.ShortenString(arg.Error(), recorder.MessageShortenLen))
			continue
		}

		shorten = append(shorten, item)
	}

	return shorten
}
