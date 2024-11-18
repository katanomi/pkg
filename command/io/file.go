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

package io

import (
	"os"
	"path"
)

// WriteFile write the content to the desired file
func WriteFile(dstFile string, content []byte, perm os.FileMode) error {
	if err := os.MkdirAll(path.Dir(dstFile), os.ModePerm); err != nil {
		return err
	}
	f, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm&os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	return err
}

// IsFile check if the given path is a file
func IsFile(file string) bool {
	s, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// IsExist check if the given path is existed
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir check if the given path is a dir
func IsDir(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return s.IsDir()
}
