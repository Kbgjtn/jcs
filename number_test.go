package jcs

import (
	"math"
	"strconv"
	"testing"
)

// see: [Appendix B: Number Serialization Samples](https://www.rfc-editor.org/rfc/rfc8785#name-number-serialization-sample)
func TestAppendNumber(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			out, err := appendNumber(nil, tc.value)
			Equals(t, tc.wantErr, err)
			Equals(t, tc.want, string(out))
		})
	}
}

func BenchmarkAppendNumber(b *testing.B) {
	samples := []float64{
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

	b.ReportAllocs()
	for _, size := range benchSizes() {
		b.Run(
			"Size="+strconv.Itoa(size),
			func(b *testing.B) {
				buf := make([]byte, 0, size*32)

				b.ResetTimer()
				for b.Loop() {
					dst := buf[:0]

					for j := 0; j < size; j++ {
						sample := samples[j%len(samples)]
						dst, _ = appendNumber(dst, sample)
					}
				}
			})
	}
}
