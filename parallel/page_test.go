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
	"testing"
)

type mockResult struct {
	Total int
	Data  []string
}

type mockRequest func(page, pageSize int) mockResult

func getMockRequestFunc(total int) mockRequest {
	return func(page, pageSize int) mockResult {
		result := mockResult{Data: make([]string, 0), Total: total}
		if total == 0 {
			return result
		}
		pageNum := total / pageSize
		for pageNum*pageSize < total {
			pageNum = pageNum + 1
		}

		if page > pageNum {
			return result
		}

		if page*pageSize <= total {
			for i := 0; i < pageSize; i++ {
				result.Data = append(result.Data, "result")
			}
		} else {
			// processing last page requests
			// for example: total 10, pageNum: 4, page: 4, pageSize 3. this page length is 1
			for i := 0; i < (total - (pageSize * (pageNum - 1))); i++ {
				result.Data = append(result.Data, "result")
			}
		}
		return result
	}
}

func TestPage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		PageRequest(context.Background(), "TestPageResult", 2, 10, PageRequestFunc{
			RequestPage: func(ctx context.Context, pageSize int, page int) (interface{}, error) {
				fmt.Printf("request -> page: %d, pagesize: %d\n", page, pageSize)
				return nil, nil
			},
			PageResult: func(items interface{}) (total int, currentPageLen int, err error) {
				return 8, 6, nil
			},
		})
	})
	t.Run("pre check mock fuck", func(t *testing.T) {
		total := 100
		mockFunc := getMockRequestFunc(total)

		preCheck1 := mockFunc(1, 10)
		// total 100, page 1, pageSize 10.
		if len(preCheck1.Data) != 10 || preCheck1.Total != total {
			t.Errorf("pre check failed, should be 10, but got %d, get total is %d", len(preCheck1.Data), preCheck1.Total)
		}

		preCheck2 := mockFunc(2, 10)
		// total 100, page 2, pageSize 10.
		if len(preCheck2.Data) != 10 || preCheck2.Total != total {
			t.Errorf("pre check failed, should be 10, but got %d, get total is %d", len(preCheck2.Data), preCheck2.Total)
		}

		total = 10
		mockFunc = getMockRequestFunc(total)
		preCheckEndPage := mockFunc(4, 3)
		// total 10, page 4, pageSize 3. should return data num is 1
		if len(preCheckEndPage.Data) != 1 || preCheckEndPage.Total != total {
			t.Errorf("pre check failed, should be 1, but got %d, get total is %d", len(preCheckEndPage.Data), preCheckEndPage.Total)
		}

		total = 100000
		mockFunc = getMockRequestFunc(total)
		preCheckBigNum := mockFunc(10000, 10)
		// total 100000, page 10000, pageSize 10. should return data num is 10000
		if len(preCheckBigNum.Data) != 10 || preCheckBigNum.Total != total {
			t.Errorf("pre check failed, should be 1, but got %d, get total is %d", len(preCheckBigNum.Data), preCheckBigNum.Total)
		}
	})
	t.Run("check result", func(t *testing.T) {
		total := 999
		mockFunc := getMockRequestFunc(total)
		overResult := make([]string, 0)

		items, _ := PageRequest(context.Background(), "CheckPageResult", 2, 10, PageRequestFunc{
			RequestPage: func(ctx context.Context, pageSize int, page int) (interface{}, error) {
				tmpResult := mockFunc(page, pageSize)
				return tmpResult, nil
			},
			PageResult: func(items interface{}) (total int, currentPageLen int, err error) {
				tmpResult := items.(mockResult)
				return tmpResult.Total, len(tmpResult.Data), nil
			},
		})

		for _, _item := range items {
			item := _item.(mockResult)
			overResult = append(overResult, item.Data...)
		}

		if len(overResult) != total {
			t.Errorf("check result failed, should be %d, got %d", total, len(overResult))
		}
	})
}

func TestPageZero(t *testing.T) {
	t.Run("check result 0", func(t *testing.T) {
		total := 0
		mockFunc := getMockRequestFunc(total)
		overResult := make([]string, 0)

		items, _ := PageRequest(context.Background(), "CheckPageResult", 2, 10, PageRequestFunc{
			RequestPage: func(ctx context.Context, pageSize int, page int) (interface{}, error) {
				tmpResult := mockFunc(page, pageSize)
				return tmpResult, nil
			},
			PageResult: func(items interface{}) (total int, currentPageLen int, err error) {
				tmpResult := items.(mockResult)
				return tmpResult.Total, len(tmpResult.Data), nil
			},
		})

		for _, _item := range items {
			item := _item.(mockResult)
			overResult = append(overResult, item.Data...)
		}

		if len(overResult) != total {
			t.Errorf("check result failed, should be %d, got %d", total, len(overResult))
		}
	})
}

func TestPageBigNum(t *testing.T) {
	t.Run("check result 0", func(t *testing.T) {
		total := 10000
		mockFunc := getMockRequestFunc(total)
		overResult := make([]string, 0)

		items, _ := PageRequest(context.Background(), "CheckPageResult", 2, 10, PageRequestFunc{
			RequestPage: func(ctx context.Context, pageSize int, page int) (interface{}, error) {
				tmpResult := mockFunc(page, pageSize)
				return tmpResult, nil
			},
			PageResult: func(items interface{}) (total int, currentPageLen int, err error) {
				tmpResult := items.(mockResult)
				return tmpResult.Total, len(tmpResult.Data), nil
			},
		})

		for _, _item := range items {
			item := _item.(mockResult)
			overResult = append(overResult, item.Data...)
		}

		if len(overResult) != total {
			t.Errorf("check result failed, should be %d, got %d", total, len(overResult))
		}
	})
}
