package jcs

import (
	"fmt"
	"testing"
	"time"
)

func TestAppendSlice(t *testing.T) {
	tests := []struct {
		name    string
		value   []any
		want    string
		wantErr error
	}{
		{"EmptySlice", []any{}, "[]", nil},
		{"SingleString", []any{"hello"}, `["hello"]`, nil},
		{"MultipleStrings", []any{"a", "b", "c"}, `["a","b","c"]`, nil},
		{"Ints", []any{1, 2, 3}, "[1,2,3]", nil}, // using fake Append
		{"Booleans", []any{true, false}, "[true,false]", nil},
		{"Times", []any{time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC), time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC)}, `["2019-01-28T07:45:10Z","2019-01-28T07:45:10.123456Z"]`, nil},
		{"MixedTypes", []any{"hi", 42, true}, "[\"hi\",42,true]", nil},
		// unsupported types
		{"ErrorUnsupportedTypes", []any{fmt.Errorf("fail")}, "", ErrUnsupportedType},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := appendSlice([]byte{}, tc.value)
			Equals(t, tc.wantErr, err)
			Equals(t, tc.want, string(out))
		})
	}
}
