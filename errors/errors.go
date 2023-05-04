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

package errors

import (
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"errors"
	"net/http"
)

// ErrNilPointer indicates nil pointer, avoid panic.
// Although unlikely just in case
var ErrNilPointer = errors.New("nil pointer")

const CredentialNotProvided metav1.StatusReason = "CredentialNotProvided"
const FileNotFound metav1.StatusReason = "FileNotFound"
const GitRevisionNotFound metav1.StatusReason = "GitRevisionNotFound"
const ToolServiceUnavailable metav1.StatusReason = "ToolServiceUnavailable"

func IsCredentialNotProvided(err error) bool {
	return k8serrors.ReasonForError(err) == CredentialNotProvided
}

func IsFileNotFound(err error) bool {
	return k8serrors.ReasonForError(err) == FileNotFound
}

func IsGitRevisionNotFound(err error) bool {
	return k8serrors.ReasonForError(err) == GitRevisionNotFound
}

func IsToolServiceUnavailable(err error) bool {
	return k8serrors.ReasonForError(err) == ToolServiceUnavailable
}

func NewCredentialNotProvided(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusUnauthorized,
			Reason:  CredentialNotProvided,
			Message: message,
		},
	}
}

func NewFileNotFound(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusNotFound,
			Reason:  FileNotFound,
			Message: message,
		},
	}
}

func NewGitRevisionNotFound(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusNotFound,
			Reason:  GitRevisionNotFound,
			Message: message,
		},
	}
}

func NewToolServiceUnavailable(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusServiceUnavailable,
			Reason:  ToolServiceUnavailable,
			Message: message,
		},
	}
}
