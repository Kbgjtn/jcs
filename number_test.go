package jcs

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

// see: [Appendix B: Number Serialization Samples](https://www.rfc-editor.org/rfc/rfc8785#name-number-serialization-sample)
func TestNumber(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		want    string
		wantErr error
	}{
		// Zero and minus zero
		{name: "Zero", value: 0.0, want: "0"},
		{name: "MinusZero", value: math.Copysign(0.0, -1.0), want: "0"},

		// NaN and Infinity must error
		{name: "NaN", value: math.NaN(), want: "", wantErr: ErrNaN},
		{name: "PosInf", value: math.Inf(1), want: "", wantErr: ErrInf},
		{name: "NegInf", value: math.Inf(-1), want: "", wantErr: ErrInf},

		// Round-to-even case
		{"RoundToEven", 1424953923781206.25, "1424953923781206.2", nil},

		{"NegativeTiny", -0.0000033333333333333333, "-0.0000033333333333333333", nil},

		// Min positive/negative subnormal
		{name: "MinPositiveNumber", value: math.SmallestNonzeroFloat64, want: "5e-324"},
		{name: "MinNegativeNumber", value: -math.SmallestNonzeroFloat64, want: "-5e-324"},

		// Max positive/negative finite
		{name: "MaxPosNumber", value: math.MaxFloat64, want: "1.7976931348623157e308"},
		{name: "MaxNegNumber", value: -math.MaxFloat64, want: "-1.7976931348623157e308"},

		// Max safe integers
		{name: "MaxPosInt", value: float64(9007199254740992), want: "9007199254740992"},
		{name: "MaxNegInt", value: float64(-9007199254740992), want: "-9007199254740992"},

		// Extended precision integer (~2**68)
		{name: "TwoPow68", value: float64(295147905179352830000), want: "295147905179352830000"},

		// Edge rounding and scientific notation
		{name: "1e+21", value: 1e+21, want: "1e21"},
		{name: "1e+23", value: 1e+23, want: "1e23"},
		{name: "0.000001", value: 0.000001, want: "0.000001"},
		{name: "9.999999999999997e-7", value: 9.999999999999997e-7, want: "9.999999999999997e-7"},
		{name: "9.999999999999997e+22", value: 9.999999999999997e+22, want: "9.999999999999997e22"},
		{name: "1.0000000000000001e+23", value: 1.0000000000000001e+23, want: "1.0000000000000001e23"},
		{name: "999999999999999700000", value: 999999999999999700000.0, want: "999999999999999700000"},
		{name: "999999999999999900000", value: 999999999999999900000.0, want: "999999999999999900000"},

		// Rounding cluster
		{name: "333333333.3333332", value: 333333333.3333332, want: "333333333.3333332"},
		{name: "333333333.3333333", value: 333333333.3333333, want: "333333333.3333333"},
		{name: "333333333.3333334", value: 333333333.3333334, want: "333333333.3333334"},
		{name: "333333333.33333343", value: 333333333.33333343, want: "333333333.33333343"},
		{name: "333333333.33333325", value: 333333333.33333325, want: "333333333.33333325"},
	}

	for _, tc := range tests {
		t.Run("append_"+tc.name, func(t *testing.T) {
			out, err := appendNumber(nil, tc.value)
			Equals(t, tc.wantErr, err)
			Equals(t, tc.want, string(out))
		})

		// WIP (dapa: enhancment encoder)
		// t.Run("encode_"+tc.name, func(t *testing.T) {
		// 	w := bytes.NewBuffer(nil)
		// 	e := NewEncoder(w)
		// 	err := e.writeNumber(tc.value)
		// 	Equals(t, tc.wantErr, err)
		// 	Equals(t, tc.want, w.String())
		// })
	}
}

func FuzzAppendNumber(f *testing.F) {
	tests := []float64{
		0,
		-0,
		1,
		-1,
		1e-9,
		-1e-9,
		1e21,
		-1e21,
		math.MaxFloat64,
		-math.MaxFloat64,
		123.456,
		-123.456,
		9007199254740992.0,
		-9007199254740992.0,
		295147905179352830000.0,
		9.999999999999997e+22,
		1e+23,
		1.0000000000000001e+23,
		999999999999999700000.0,
		999999999999999900000.0,
		1e+21,
		9.999999999999997e-7,
		0.000001,
		333333333.3333332,
		333333333.33333325,
		333333333.3333333,
		333333333.3333334,
		333333333.33333343,
		-0.0000033333333333333333,
		1424953923781206.25,
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		got, err := appendNumber(nil, v)
		if err != nil {
			switch {
			case errors.Is(err, ErrInf),
				errors.Is(err, ErrNumberOOR),
				errors.Is(err, ErrNaN):

				return
			}

			t.Fatal("unexpected error:", err)
			return
		}

		want, _ := json.Marshal(v)
		want = bytes.ReplaceAll(want, []byte("e+"), []byte("e"))
		want = bytes.ReplaceAll(want, []byte("E+"), []byte("E"))

		if !bytes.Equal(want, got) {
			t.Fatalf("input=%v\nwant=%q (len=%d)\ngot =%q (len=%d)\nwant hex=%x\ngot hex=%x", v, want, len(want), got, len(got), want, got)
			return
		}
	})
}

var NumbersSample = []float64{
	0.0,
	math.Copysign(0.0, -1.0), // -0
	math.SmallestNonzeroFloat64,
	-math.SmallestNonzeroFloat64,
	math.MaxFloat64,
	-math.MaxFloat64,
	9007199254740992.0,
	-9007199254740992.0,
	295147905179352830000.0,
	9.999999999999997e+22,
	1e+23,
	1.0000000000000001e+23,
	999999999999999700000.0,
	999999999999999900000.0,
	1e+21,
	9.999999999999997e-7,
	0.000001,
	333333333.3333332,
	333333333.33333325,
	333333333.3333333,
	333333333.3333334,
	333333333.33333343,
	-0.0000033333333333333333,
	1424953923781206.25,
}

func BenchmarkAppendNumber_Throughput(b *testing.B) {
	for _, size := range benchSizes() {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			buf := make([]byte, 0, size*AvgBytesPerNumber)

			b.ReportAllocs()
			b.ResetTimer()
			b.SetBytes(int64(size) * int64(AvgBytesPerNumber))

			for b.Loop() {
				dst := buf[:0]
				for i := 0; i < size; i++ {
					sample := NumbersSample[i%len(NumbersSample)]
					dst, _ = appendNumber(dst, sample)
				}
			}
		})
	}
}

func BenchmarkAppendNumber(b *testing.B) {
	avgBytesPerNumber := 16
	dst := make([]byte, 0, avgBytesPerNumber)
	var sample float64

	b.ReportAllocs()
	for _, size := range benchSizes() {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			b.SetBytes(int64(avgBytesPerNumber))
			for i := 0; i < b.N; i++ {
				dst = dst[:0]
				sample = NumbersSample[i%len(NumbersSample)]
				dst, _ = appendNumber(dst, sample)
			}
		})
	}
}
