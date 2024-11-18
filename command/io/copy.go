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
	"io"
	"io/ioutil"
	"os"
	"path"
)

// copyFile copies a files from src to dst.
func copyFile(src, dst string, srcInfo os.FileInfo, config *Config) error {
	for _, filter := range config.FileFilter {
		if filter(src, dst) {
			return nil
		}
	}
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	if IsDir(dst) {
		dst = path.Join(dst, path.Base(src))
	} else {
		if err = os.MkdirAll(path.Dir(dst), os.ModePerm); err != nil {
			return err
		}
	}

	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, srcInfo.Mode()&os.ModePerm)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

// copyDir copies a dir from src to dst
func copyDir(src, dst string, srcInfo os.FileInfo, config *Config) error {
	var err error
	var fds []os.FileInfo

	if err = os.MkdirAll(dst, srcInfo.Mode()&os.ModePerm); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if err = switchboard(srcFp, dstFp, fd, config); err != nil {
			return err
		}
	}
	return nil
}

// Copy support copy file to file or dir to dir
func Copy(src, dst string, opts ...CopyOption) error {
	config := &Config{}
	for _, option := range opts {
		option(config)
	}
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return switchboard(src, dst, info, config)
}

func copyLink(src, dest string, config *Config) (err error) {
	var orig string
	orig, err = os.Readlink(src)
	if err != nil {
		return err
	}
	if !path.IsAbs(orig) {
		orig = path.Join(path.Dir(src), orig)
	}
	var info os.FileInfo
	info, err = os.Lstat(orig)
	if err != nil {
		return err
	}
	return switchboard(orig, dest, info, config)
}

func switchboard(src, dest string, info os.FileInfo, config *Config) (err error) {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		err = copyLink(src, dest, config)
	case info.IsDir():
		err = copyDir(src, dest, info, config)
	default:
		err = copyFile(src, dest, info, config)
	}

	return err
}
