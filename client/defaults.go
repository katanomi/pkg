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

package client

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	DefaultTimeout         = 10 * time.Second
	DefaultQPS     float32 = 50.0
	DefaultBurst           = 60
)

func NewHTTPClient() *http.Client {
	var timeout int64
	timeoutStr := os.Getenv("HTTP_CLIENT_TIMEOUT")
	timeout, err := strconv.ParseInt(timeoutStr, 10, 64)
	if len(timeoutStr) == 0 || err != nil {
		timeout = 30
	}

	client := &http.Client{
		Transport: GetDefaultTransport(),
		Timeout:   time.Duration(timeout) * time.Second,
	}

	return client
}

func GetDefaultTransport() http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
