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

package v1alpha1

type CodeScanReportSyncReason string

const (
	// Code scan task succeeded
	CodeScanReportSyncSucceededReason CodeScanReportSyncReason = SucceededReason
	// The code scan task was canceled
	CodeScanTaskCancelled CodeScanReportSyncReason = "CodeScanTaskCancelled"
	// Code scan task is waiting to start
	CodeScanTaskPending CodeScanReportSyncReason = "CodeScanTaskPending"
	// ode scan task in process
	CodeScanTaskRunning CodeScanReportSyncReason = "CodeScanTaskRunning"
	// Code scan task execution failed
	CodeScanTaskFailedReason CodeScanReportSyncReason = "CodeScanTaskFailed"
	// Error in sync report
	CodeScanReportSyncErrorReason CodeScanReportSyncReason = "CodeScanSyncReportError"
)
