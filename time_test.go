package jcs

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestAppendTime(t *testing.T) {
	tests := []struct {
		name   string
		in     time.Time
		exp    string
		expErr error
	}{
		{
			name:   "2019-01-28T07:45:10Z",
			in:     time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC),
			exp:    "2019-01-28T07:45:10Z",
			expErr: nil,
		},
		{
			name:   "2019-01-28T07:45:10.123456Z",
			in:     time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC),
			exp:    "2019-01-28T07:45:10Z",
			expErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := []byte{}
			buf = appendTime(buf, tc.in)
		})
	}
}

func BenchmarkAppendTime(b *testing.B) {
	b.ReportAllocs()

	base := time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, size := range benchSizes() {
		buf := make([]byte, 0, size*64)

		b.Run(
			"Size="+strconv.Itoa(size),
			func(b *testing.B) {
				dst := buf[:0]

				b.ResetTimer()
				for b.Loop() {
					randomOffset := time.Duration(rng.Int63n(int64(time.Hour*24*365*2))) - time.Duration(time.Hour*24*365)
					randomTime := base.Add(randomOffset)

					_ = appendTime(dst, randomTime)
				}
			},
		)
	}
}
