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
	"errors"
	"testing"

	. "github.com/onsi/gomega"
)

type testStruct struct {
	name string
}

func mockPaging(total int) FetchFunc {
	return func(_ context.Context, pageSize, page int) (int, interface{}, error) {
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > total {
			end = total
		}
		result := make([]int, 0)
		for i := start; i < end; i++ {
			result = append(result, i+1)
		}
		return total, result, nil
	}
}

func Test_mockPaging(t *testing.T) {
	ctx := context.Background()
	g := NewGomegaWithT(t)
	tests := []struct {
		total    int
		page     int
		pageSize int
		list     interface{}
	}{
		{
			total:    10,
			page:     1,
			pageSize: 5,
			list:     []interface{}{1, 2, 3, 4, 5},
		},
		{
			total:    10,
			page:     2,
			pageSize: 5,
			list:     []interface{}{6, 7, 8, 9, 10},
		},
		{
			total:    8,
			page:     1,
			pageSize: 3,
			list:     []interface{}{1, 2, 3},
		},
		{
			total:    8,
			page:     2,
			pageSize: 3,
			list:     []interface{}{4, 5, 6},
		},
		{
			total:    8,
			page:     3,
			pageSize: 3,
			list:     []interface{}{7, 8},
		},
	}

	for _, tt := range tests {
		gotTotal, gotList, _ := mockPaging(tt.total)(ctx, tt.pageSize, tt.page)
		g.Expect(gotTotal).To(Equal(tt.total))
		g.Expect(expandSlice(gotList)).To(Equal(tt.list))
	}
}

func Test_calcPageCount(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		total         int
		pageSize      int
		wantPageCount int
	}{
		{
			total:         0,
			pageSize:      0,
			wantPageCount: 0,
		},
		{
			total:         -1,
			pageSize:      -1,
			wantPageCount: 0,
		},
		{
			total:         1,
			pageSize:      10,
			wantPageCount: 1,
		},
		{
			total:         10,
			pageSize:      10,
			wantPageCount: 1,
		},
		{
			total:         21,
			pageSize:      10,
			wantPageCount: 3,
		},
	}
	for _, tt := range tests {
		g.Expect(calcPageCount(tt.total, tt.pageSize)).To(Equal(tt.wantPageCount))
	}
}

func Test_getCount(t *testing.T) {
	g := NewGomegaWithT(t)
	tests := []struct {
		data      interface{}
		wantCount int
		wantOk    bool
	}{
		{
			data:      1,
			wantCount: 0,
			wantOk:    false,
		},
		{
			data:      "abc",
			wantCount: 0,
			wantOk:    false,
		},
		{
			data:      testStruct{},
			wantCount: 0,
			wantOk:    false,
		},
		{
			data:      []interface{}(nil),
			wantCount: 0,
			wantOk:    true,
		},
		{
			data:      []int{},
			wantCount: 0,
			wantOk:    true,
		},
		{
			data:      []int{1, 2, 3, 4},
			wantCount: 4,
			wantOk:    true,
		},
		{
			data:      []testStruct{{}, {}, {}},
			wantCount: 3,
			wantOk:    true,
		},
	}
	for _, tt := range tests {
		gotCount, gotOk := getCount(tt.data)
		g.Expect(gotCount).To(Equal(tt.wantCount))
		g.Expect(gotOk).To(Equal(tt.wantOk))
	}
}

func Test_expandSlice(t *testing.T) {
	t1, t2, t3 := testStruct{}, testStruct{}, testStruct{}
	g := NewGomegaWithT(t)
	tests := []struct {
		data        interface{}
		wantResults []interface{}
	}{
		{
			data:        []interface{}(nil),
			wantResults: []interface{}{},
		},
		{
			data:        []string{"a", "b"},
			wantResults: []interface{}{"a", "b"},
		},
		{
			data:        []int{1, 2},
			wantResults: []interface{}{1, 2},
		},
		{
			data:        []testStruct{t1, t2, t3},
			wantResults: []interface{}{t1, t2, t3},
		},
		{
			data:        []*testStruct{&t1, &t2, &t3},
			wantResults: []interface{}{&t1, &t2, &t3},
		},
	}
	for _, tt := range tests {
		g.Expect(expandSlice(tt.data)).To(Equal(tt.wantResults))
	}
}

func TestFetchAllPage_success(t *testing.T) {
	ctx := context.Background()
	g := NewGomegaWithT(t)
	total := 23
	wantList := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	paging := mockPaging(total)
	for pageSize := 1; pageSize < total; pageSize++ {
		got, err := FetchAllPage(ctx, 2, pageSize, paging)
		g.Expect(got).To(Equal(wantList))
		g.Expect(err).Should(Succeed())
	}
}

func TestFetchAllPage_fail(t *testing.T) {
	ctx := context.Background()
	g := NewGomegaWithT(t)
	testErr := errors.New("test error")

	genFetchFun := func(total, errPage int) FetchFunc {
		return func(ctx context.Context, pageSize, page int) (int, interface{}, error) {
			if page == errPage {
				return 0, nil, testErr
			}
			return mockPaging(total)(ctx, pageSize, page)
		}
	}

	tests := []struct {
		pageSize int
		total    int
		errPage  int
		wantErr  bool
	}{
		{
			pageSize: 5,
			total:    11,
			errPage:  1,
			wantErr:  true,
		},
		{
			pageSize: 5,
			total:    11,
			errPage:  2,
			wantErr:  true,
		},
		{
			pageSize: 5,
			total:    11,
			errPage:  3,
			wantErr:  true,
		},
		{
			pageSize: 5,
			total:    11,
			errPage:  4,
			wantErr:  false,
		},
		{
			pageSize: 5,
			total:    111,
			errPage:  20,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		_, err := FetchAllPage(ctx, 2, tt.pageSize, genFetchFun(tt.total, tt.errPage))
		g.Expect(err != nil).To(Equal(tt.wantErr))
	}
}
