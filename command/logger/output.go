/*
Copyright 2022 The AlaudaDevops Authors.

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

package logger

import (
	"io"
	"os"
)

// OutputWriter describe an output writer
type OutputWriter interface {
	Output(s string) error
}

// NewOutput construct an OutputWriter
func NewOutput(writer io.Writer) OutputWriter {
	if writer == nil {
		return &output{writer: os.Stdout}
	}
	return &output{writer: writer}
}

type output struct {
	writer io.Writer
}

// Output aliases for WriteString
func (o *output) Output(s string) error {
	return o.WriteString(s)
}

// WriteString appends the contents of s to o's writer.
func (o *output) WriteString(s string) (err error) {
	return o.Write([]byte(s))
}

// Write appends the byte c to o's writer.
func (o *output) Write(bytes []byte) (err error) {
	_, err = o.writer.Write(bytes)
	return err
}
