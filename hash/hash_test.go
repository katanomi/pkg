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

package hash

import (
	"fmt"
	"hash/adler32"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type A struct {
	x int
	y string
}

type B struct {
	x []int
	y map[string]bool
}

type C struct {
	x int
	y string
}

func (c C) String() string {
	return fmt.Sprintf("%d:%s", c.x, c.y)
}

func TestDeepHashObject(t *testing.T) {
	successCases := []func() interface{}{
		func() interface{} { return 8675309 },
		func() interface{} { return "Jenny, I got your number" },
		func() interface{} { return []string{"eight", "six", "seven"} },
		func() interface{} { return [...]int{5, 3, 0, 9} },
		func() interface{} { return map[int]string{8: "8", 6: "6", 7: "7"} },
		func() interface{} { return map[string]int{"5": 5, "3": 3, "0": 0, "9": 9} },
		func() interface{} { return A{867, "5309"} },
		func() interface{} { return &A{867, "5309"} },
		func() interface{} {
			return B{[]int{8, 6, 7}, map[string]bool{"5": true, "3": true, "0": true, "9": true}}
		},
		func() interface{} { return map[A]bool{{8675309, "Jenny"}: true, {9765683, "!Jenny"}: false} },
		func() interface{} { return map[C]bool{{8675309, "Jenny"}: true, {9765683, "!Jenny"}: false} },
		func() interface{} { return map[*A]bool{{8675309, "Jenny"}: true, {9765683, "!Jenny"}: false} },
		func() interface{} { return map[*C]bool{{8675309, "Jenny"}: true, {9765683, "!Jenny"}: false} },
	}

	for _, tc := range successCases {
		hasher1 := adler32.New()
		DeepHashObject(hasher1, tc())
		hash1 := hasher1.Sum32()
		DeepHashObject(hasher1, tc())
		hash2 := hasher1.Sum32()
		if hash1 != hash2 {
			t.Fatalf("hash of the same object (%q) produced different results: %d vs %d", toString(tc()), hash1, hash2)
		}
		for i := 0; i < 100; i++ {
			hasher2 := adler32.New()

			DeepHashObject(hasher1, tc())
			hash1a := hasher1.Sum32()
			DeepHashObject(hasher2, tc())
			hash2a := hasher2.Sum32()

			if hash1a != hash1 {
				t.Errorf("repeated hash of the same object (%q) produced different results: %d vs %d", toString(tc()), hash1, hash1a)
			}
			if hash2a != hash2 {
				t.Errorf("repeated hash of the same object (%q) produced different results: %d vs %d", toString(tc()), hash2, hash2a)
			}
			if hash1a != hash2a {
				t.Errorf("hash of the same object produced (%q) different results: %d vs %d", toString(tc()), hash1a, hash2a)
			}
		}
	}
}

func toString(obj interface{}) string {
	return spew.Sprintf("%#v", obj)
}

type wheel struct {
	radius uint32
}

type unicycle struct {
	primaryWheel   *wheel
	licencePlateID string
	tags           map[string]string
}

func TestDeepObjectPointer(t *testing.T) {
	// Arrange
	wheel1 := wheel{radius: 17}
	wheel2 := wheel{radius: 22}
	wheel3 := wheel{radius: 17}

	myUni1 := unicycle{licencePlateID: "blah", primaryWheel: &wheel1, tags: map[string]string{"color": "blue", "name": "john"}}
	myUni2 := unicycle{licencePlateID: "blah", primaryWheel: &wheel2, tags: map[string]string{"color": "blue", "name": "john"}}
	myUni3 := unicycle{licencePlateID: "blah", primaryWheel: &wheel3, tags: map[string]string{"color": "blue", "name": "john"}}

	// Run it more than once to verify determinism of hasher.
	for i := 0; i < 100; i++ {
		hasher1 := adler32.New()
		hasher2 := adler32.New()
		hasher3 := adler32.New()
		// Act
		DeepHashObject(hasher1, myUni1)
		hash1 := hasher1.Sum32()
		DeepHashObject(hasher1, myUni1)
		hash1a := hasher1.Sum32()
		DeepHashObject(hasher2, myUni2)
		hash2 := hasher2.Sum32()
		DeepHashObject(hasher3, myUni3)
		hash3 := hasher3.Sum32()

		// Assert
		if hash1 != hash1a {
			t.Errorf("repeated hash of the same object produced different results: %d vs %d", hash1, hash1a)
		}

		if hash1 == hash2 {
			t.Errorf("hash1 (%d) and hash2(%d) must be different because they have different values for wheel size", hash1, hash2)
		}

		if hash1 != hash3 {
			t.Errorf("hash1 (%d) and hash3(%d) must be the same because although they point to different objects, they have the same values for wheel size", hash1, hash3)
		}
	}
}

func TestComputeHash(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{
				obj: nil,
			},
			want: "99c9b8c64",
		},
		{
			name: "empty string",
			args: args{
				obj: "",
			},
			want: "8857b86c7",
		},
		{
			name: "number",
			args: args{
				obj: 1234,
			},
			want: "6ff6c884f5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeHash(tt.args.obj); got != tt.want {
				t.Errorf("ComputeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashSHA256(t *testing.T) {
	type args struct {
		secretKey string
		value     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty secretKey and empty value",
			args: args{
				secretKey: "",
				value:     []byte{},
			},
			want:    "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
			wantErr: false,
		},
		{
			name: "empty secretKey",
			args: args{
				secretKey: "",
				value:     []byte{'a', 'b', 'c', 'd'},
			},
			want:    "527ff4c28c22a090fe39908139363e81b8fb10d0695a135518006abfa21cf5a2",
			wantErr: false,
		},
		{
			name: "empty value",
			args: args{
				secretKey: "abcd",
				value:     []byte{},
			},
			want:    "2722000cbc34892ac64a8fb9ef2b50fc824ea1984cb81e50d687648f2e88f724",
			wantErr: false,
		},
		{
			name: "non-empty secretKey and value",
			args: args{
				secretKey: "abcd",
				value:     []byte{'a', 'b', 'c', 'd'},
			},
			want:    "e1a20dce13e4953e3d50e7f6651a0ce862a655fc84c35352447eff99a5a02852",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashSHA256(tt.args.secretKey, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashSHA256() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}
