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
	"reflect"

	"knative.dev/pkg/logging"
)

// FetchFunc fetch a page of data
// The return value `list` must be a slice or panic will occur
type FetchFunc func(ctx context.Context, pageSize, page int) (total int, list interface{}, err error)

// FetchAllPage get all data concurrently
func FetchAllPage(ctx context.Context, concurrency int, pageSize int, fetchFun FetchFunc) ([]interface{}, error) {
	logger := logging.FromContext(ctx)
	total, data, err := fetchFun(ctx, pageSize, 1)
	if err != nil {
		return nil, err
	}

	count, ok := getCount(data)
	if !ok {
		err = fmt.Errorf("fetch function expect a slice but got %s", reflect.TypeOf(data))
		logger.Error(err.Error())
		return nil, err
	}
	if count >= total || count < pageSize {
		return expandSlice(data), nil
	}

	totalPages := calcPageCount(total, pageSize)
	p := P(logger, "PageRequest").FailFast().SetConcurrent(concurrency).Context(ctx)
	for i := 2; i <= totalPages; i++ {
		page := i
		p.Add(func() (interface{}, error) {
			_, _list, _err := fetchFun(ctx, pageSize, page)
			return _list, _err
		})
	}
	results, err := p.Do().Wait()
	if err != nil {
		return nil, err
	}

	list := expandSlice(data)
	for _, result := range results {
		list = append(list, expandSlice(result)...)
	}

	return list, nil
}

// expandSlice expand an interface{} to []interface{}
// NOTE: `data` param must be a slice
func expandSlice(data interface{}) (results []interface{}) {
	v := reflect.ValueOf(data)
	results = make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		results = append(results, v.Index(i).Interface())
	}
	return results
}

func getCount(data interface{}) (count int, ok bool) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return 0, false
	}
	return v.Len(), true
}

func calcPageCount(total, pageSize int) (count int) {
	if total <= 0 || pageSize <= 0 {
		return 0
	}
	if total <= pageSize {
		return 1
	}
	count = total / pageSize
	if total%pageSize != 0 {
		count++
	}
	return count
}
