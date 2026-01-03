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
			name: "zero",
			in:   time.Time{},
			want: `"0001-01-01T00:00:00Z"`,
		},
		{
			name: "epoch",
			in:   time.Unix(0, 0).UTC(),
			want: `"1970-01-01T00:00:00Z"`,
		},
		{
			name: "far future",
			in:   time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			want: `"9999-12-31T23:59:59.999999999Z"`,
		},
		{
			name: "far past",
			in:   time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			want: `"0001-01-01T00:00:00Z"`,
		},
		{
			name: "end of February",
			in:   time.Date(2019, 2, 28, 23, 59, 59, 0, time.UTC),
			want: `"2019-02-28T23:59:59Z"`,
		},
		{
			name: "leap day with ns",
			in:   time.Date(2020, 2, 29, 12, 0, 0, 123456789, time.UTC),
			want: `"2020-02-29T12:00:00.123456789Z"`,
		},
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

		want, _ := json.Marshal(tc.in)
		t.Log("want:", string(want))

		t.Run("append_"+tc.name, func(t *testing.T) {
			got := appendTime([]byte{}, tc.in)
			Equals(t, tc.want, string(got))
		})

		// WIP (dapa: enhancment encoder)
		// t.Run("write_"+tc.name, func(t *testing.T) {
		// 	w := bytes.NewBuffer(nil)
		// 	e := NewEncoder(w)
		// 	err := e.writeTime(tc.in)
		// 	Equals(t, nil, err)
		// 	Equals(t, tc.want, w.String())
		// })
	}
}

func FuzzAppendTime(f *testing.F) {
	tests := []int64{
		time.Time{}.UnixNano(),     // zero
		time.Unix(0, 0).UnixNano(), // epoch
		time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC).UnixNano(), // far future
		time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano(),                 // far past
		time.Date(2019, 2, 28, 23, 59, 59, 0, time.UTC).UnixNano(),          // end of Feb (non-leap)
		time.Date(2020, 2, 29, 12, 0, 0, 123456789, time.UTC).UnixNano(),    // leap day with ns
		time.Date(2019, 1, 28, 7, 45, 10, 0, time.UTC).UnixNano(),
		time.Date(2019, 1, 28, 7, 45, 10, 123456000, time.UTC).UnixNano(),
	}

	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, nanos int64) {
		s := time.Unix(0, nanos).UTC()
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
