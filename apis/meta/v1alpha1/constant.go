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

type ReportSyncReason string

const (
	// SkippedReason a task/stage/step was skipped
	SkippedReason = "Skipped"
	// CancelledReason a task/stage/step was cancelled
	CancelledReason = "Cancelled"
	// PendingReason resource is in a pending state waiting for some condition
	PendingReason = "Pending"
	// ApprovedReason a request/approval was approved
	ApprovedReason = "Approved"
	// DeniedReason a request/approval was denied
	DeniedReason = "Denied"
	// NotRequired a request/approval is not required and therefore was "approved"
	NotRequiredReason = "NotRequired"
	// RunningReason entered in a running state, a quality gate check is executing
	RunningReason = "Running"
	// CreatedReason created the object
	CreatedReason = "Created"
	// FailedReason execution/tests/quality gates failed
	FailedReason = "Failed"
	// SucceededReason execution/tests/quality gates succeeded
	SucceededReason = "Succeeded"
	// TimeoutReason execution/tests/quality gates timed out
	TimeoutReason = "Timeout"
	// QueuedReason is in a queued state
	QueuedReason = "Queued"
	// ValidationError some validation error occurred
	ValidationError = "ValidationError"
	// ErrorReason some error occurred
	ErrorReason = "Error"
)

const (
	// Sync task report succeeded
	ReportSyncSucceededReason ReportSyncReason = SucceededReason
	// Report is syncing
	ReportSyncingReason ReportSyncReason = "ReportSyncing"
	// Failed to sync report
	ReportSyncFailedReason ReportSyncReason = "ReportSyncFailed"
)
