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

package config

// WatchFunc is the callback function for config watcher
type WatchFunc func(*Config)

// Watcher describes the interface for config watcher
type Watcher interface {
	Watch(config *Config)
}

// NewConfigWatcher constructs a new config watcher
// please note that the callback should be processed quickly
// and should not block
func NewConfigWatcher(cb WatchFunc) Watcher {
	c := &configWatcher{
		inputs:   make(chan *Config),
		callback: cb,
	}
	go c.run()
	return c
}

// configWatcher is the implementation of config watcher
// it will call the callback function serially when a new config is received
type configWatcher struct {
	inputs   chan *Config
	callback WatchFunc
}

func (p configWatcher) run() {
	for data := range p.inputs {
		if p.callback != nil {
			p.callback(data)
		}
	}
}

// Watch implements the Watcher interface
func (p configWatcher) Watch(config *Config) {
	p.inputs <- config
}
