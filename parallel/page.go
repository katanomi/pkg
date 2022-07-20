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

package parallel

import (
	"context"
	"fmt"
	"time"

	"k8s.io/utils/trace"
	"knative.dev/pkg/logging"
)

// PageRequestFunc is a tool for concurrent processing of pagination
type PageRequestFunc struct {
	// RequestPage for concurrent request paging
	RequestPage func(ctx context.Context, pageSize int, page int) (interface{}, error)
	// PageResult for get paging information
	PageResult func(items interface{}) (total int, currentPageLen int, err error)
}

// PageRequest is concurrent request paging
func PageRequest(ctx context.Context, logName string, concurrency int, pageSize int, f PageRequestFunc) ([]interface{}, error) {
	log := trace.New("PageRequest", trace.Field{Key: "name", Value: logName})
	logger := logging.FromContext(ctx)

	defer func() {
		log.LogIfLong(3 * time.Second)
	}()

	items, err := f.RequestPage(ctx, pageSize, 1)
	if err != nil {
		return nil, err
	}
	log.Step("requested page 1")
	total, firstPageLen, err := f.PageResult(items)
	if err != nil {
		return nil, err
	}
	if firstPageLen < pageSize {
		return []interface{}{items}, nil
	}

	if total == firstPageLen {
		return []interface{}{items}, nil
	}

	var request = func(i int) func() (interface{}, error) {
		return func() (interface{}, error) {
			items, err := f.RequestPage(ctx, pageSize, i)
			log.Step(fmt.Sprintf("requested page %d", i))
			return items, err
		}
	}

	totalPage := total / pageSize
	if total%pageSize != 0 {
		totalPage = totalPage + 1
	}

	if totalPage-1 < concurrency { // first page we have requested, so skip first page
		concurrency = totalPage - 1
	}

	p := P(logger, "PageRequest").FailFast().SetConcurrent(concurrency).Context(ctx)
	for i := 2; i <= totalPage; i++ {
		p.Add(request(i))
	}

	results, err := p.Do().Wait()
	if err != nil {
		return nil, err
	}

	return append([]interface{}{items}, results...), nil

}
