/*
Copyright 2023 The Katanomi Authors.

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

	"net/http"
)

const (
	// StatusReasonCredentialNotProvided indicate that credential not provided
	StatusReasonCredentialNotProvided metav1.StatusReason = "CredentialNotProvided" //nolint:gosec
	// StatusReasonUnauthorized indicate that the credential provided is invalid
	StatusReasonUnauthorized metav1.StatusReason = "Unauthorized"

	// StatusReasonFileNotFound indicate that requested file not found
	StatusReasonFileNotFound metav1.StatusReason = "FileNotFound"

	// StatusReasonGitRevisionNotFound indicate that requested git revision not found
	StatusReasonGitRevisionNotFound metav1.StatusReason = "GitRevisionNotFound"

	// StatusReasonToolServiceUnavailable indicate that requested tool service unavailable
	StatusReasonToolServiceUnavailable metav1.StatusReason = "ToolServiceUnavailable"

	// StatusReasonStorageClassNotFound indicate that default storage class not found
	StatusReasonStorageClassNotFound metav1.StatusReason = "DefaultStorageClassNotFound"
)

// NewCredentialNotProvided init a CredentialNotProvided k8s api error
func NewCredentialNotProvided(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusUnauthorized,
			Reason:  StatusReasonCredentialNotProvided,
			Message: message,
		},
	}
}

// NewFileNotFound init a FileNotFound k8s api error
func NewFileNotFound(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusNotFound,
			Reason:  StatusReasonFileNotFound,
			Message: message,
		},
	}
}

// NewGitRevisionNotFound init a GitRevisionNotFound k8s api error
func NewGitRevisionNotFound(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusNotFound,
			Reason:  StatusReasonGitRevisionNotFound,
			Message: message,
		},
	}
}

// NewDefaultStorageClassNotFound init a GitRevisionNotFound k8s api error
func NewDefaultStorageClassNotFound(message string) error {
	return &k8serrors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    http.StatusNotFound,
			Reason:  StatusReasonStorageClassNotFound,
			Message: message,
		},
	}
}

// IsCredentialNotProvided judge if the error is CredentialNotProvided
func IsCredentialNotProvided(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonCredentialNotProvided
}

// IsUnauthorized judge if the error is Unauthorized
func IsUnauthorized(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonUnauthorized
}

// IsFileNotFound judge if the error is FileNotFound
func IsFileNotFound(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonFileNotFound
}

// IsGitRevisionNotFound judge if the error is GitRevisionNotFound
func IsGitRevisionNotFound(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonGitRevisionNotFound
}

// IsToolServiceUnavailable judge if the error is ToolServiceUnavailable
func IsToolServiceUnavailable(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonToolServiceUnavailable
}

// IsDefaultStorageClassNotFound judge if the error is DefaultStorageClassNotFound
func IsDefaultStorageClassNotFound(err error) bool {
	return k8serrors.ReasonForError(err) == StatusReasonStorageClassNotFound
}

// Reason return reason for the error if it is a status reason error
func Reason(err error) (exist bool, reason metav1.StatusReason) {
	if IsCredentialNotProvided(err) {
		return true, StatusReasonCredentialNotProvided
	}
	if IsFileNotFound(err) {
		return true, StatusReasonFileNotFound
	}
	if IsGitRevisionNotFound(err) {
		return true, StatusReasonGitRevisionNotFound
	}
	if IsToolServiceUnavailable(err) {
		return true, StatusReasonToolServiceUnavailable
	}
	if IsDefaultStorageClassNotFound(err) {
		return true, StatusReasonStorageClassNotFound
	}
	if IsUnauthorized(err) {
		return true, StatusReasonUnauthorized
	}
	return false, metav1.StatusReasonUnknown
}
