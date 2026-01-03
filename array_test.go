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

func BenchmarkAppendSlice(b *testing.B) {
	avgBytesPerItem := 64
	samples := []any{
		"hi",
		"foobar",
		"e9afe14c-8421-4b71-af02-3f9d0238f1d8",
		42,
		float64(0.99),
		uint(1000),
		true,
		false,
		nil,
		int32(32332),
		time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC),
	}

	buf := make([]byte, 0)
	var sample []any
	for _, size := range sliceSize {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			sample = make([]any, size)
			for i := 0; i < size; i++ {
				sample[i] = samples[i%len(samples)]
			}

			est := int64(size) * int64(avgBytesPerItem)

			b.ReportAllocs()
			b.ResetTimer()
			b.SetBytes(est)
			for b.Loop() {
				buf = buf[:0]
				buf, _ = appendSlice(buf, sample)
			}
		})
	}
}
