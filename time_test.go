package jcs

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want string
	}{
		{
			name: "2019-01-28T07:45:10Z",
			in:   time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC),
			want: `"2019-01-28T07:45:10Z"`,
		},
		{
			name: "2019-01-28T07:45:10.123456Z",
			in:   time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC),
			want: `"2019-01-28T07:45:10.123456Z"`,
		},
	}

	for _, tc := range tests {
		t.Run("append_"+tc.name, func(t *testing.T) {
			got := appendTime([]byte{}, tc.in)
			Equals(t, tc.want, string(got))
		})

		t.Run("write_"+tc.name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)
			e := NewEncoder(w)
			err := e.writeTime(tc.in)
			Equals(t, nil, err)
			Equals(t, tc.want, w.String())
		})
	}
}

func FuzzAppendTime(f *testing.F) {
	tests := []int64{
		time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC).UnixNano(),
		time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC).UnixNano(),
	}

	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, nanos int64) {
		s := time.Unix(0, nanos)
		got := appendTime(nil, s)
		want, _ := json.Marshal(s)
		if !bytes.Equal(want, got) {
			t.Fatalf("input=%q\nwant=%q (len=%d)\ngot =%q (len=%d)\nwant hex=%x\ngot hex=%x", s, want, len(want), got, len(got), want, got)
			return
		}
	})
}

func BenchmarkAppendTime(b *testing.B) {
	b.ReportAllocs()

	base := time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, size := range benchSizes() {
		buf := make([]byte, 0, size*64)

		b.Run("Size="+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				buf = buf[:0]
				randomOffset := time.Duration(rng.Int63n(int64(time.Hour*24*365*2))) - time.Duration(time.Hour*24*365)
				randomTime := base.Add(randomOffset)

				buf = appendTime(buf, randomTime)
			}
		},
		)
	}
}
