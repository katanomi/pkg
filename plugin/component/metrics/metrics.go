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

package metrics

import (
	"strconv"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus"
)

// counter prometheus counter, for counting request count
var counter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Subsystem: "plugin api",
		Name:      "request_total",
		Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"code", "method", "url"},
)

// histogram prometheus histogram, for recording request latency
var histogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Subsystem: "plugin api",
		Name:      "request_latency",
		Help:      "The HTTP request latencies in seconds.",
	},
	[]string{"code", "method", "url"},
)

// Filter prometheus metric filter for go restful
func Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()

	chain.ProcessFilter(req, resp)

	status := strconv.Itoa(resp.StatusCode())
	latency := float64(time.Since(start)) / float64(time.Second)
	method := req.Request.Method
	url := req.Request.URL.String()

	histogram.WithLabelValues(status, method, url).Observe(latency)
	counter.WithLabelValues(status, method, url).Inc()
}
