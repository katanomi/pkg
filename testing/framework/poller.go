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

package framework

import "time"

// Poller describe the configuration of Poller
// Interval: Indicates the inspection interval
// Timeout: Indicates the maximum time the check will last
type Poller struct {
	Interval time.Duration
	Timeout  time.Duration
}

// Settings get the configuration
// Default values will be used if no custom settings are made
func (p Poller) Settings() (interval, timeout time.Duration) {
	interval = p.Interval
	timeout = p.Timeout
	if p.Interval == 0 {
		interval = 200 * time.Millisecond
	}

	if p.Timeout == 0 {
		timeout = 5 * time.Second
	}

	if interval > timeout {
		timeout = interval
	}

	return
}
