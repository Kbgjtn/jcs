package jcs

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr error
	}{
		{"Empty", "", `""`, nil},
		{"SimpleASCII", "simple ascii", `"simple ascii"`, nil},
		{"QuoteEscape", `quote"slash\test`, `"quote\"slash\\test"`, nil},
		{"ControlChars", "control:\b\t\n\f\r", `"control:\b\t\n\f\r"`, nil},
		{"LowUnicode", "low\u0001\u0002\u001f", `"low\u0001\u0002\u001f"`, nil},
		{"NonASCII", "こんにちは世界", `"こんにちは世界"`, nil},
		{"Emoji", "emoji 😀😅🚀", `"emoji 😀😅🚀"`, nil},
		{"Mixed", "mixed ascii 日本語 😀 \n \" \\", `"mixed ascii 日本語 😀 \n \" \\"`, nil},
		{"U+FFFD valid", string([]byte{0xEF, 0xBF, 0xBD}), string("\"\uFFFD\""), nil}, // replacement character
		{"Invalid single-byte sequences", string([]byte{0xFF}), "", ErrInvalidUTF8},   // invalid UTF-8 byte
		{"truncated 3-byte sequence", string([]byte{0xE2, 0x82}), "", ErrInvalidUTF8},
		{"SurrogateHigh", string([]byte{0xED, 0xA0, 0x80}), "", ErrInvalidUTF8},
		{"SurrogateLow", string([]byte{0xED, 0xBF, 0xBF}), "", ErrInvalidUTF8},
		{"OverlongEncoding", string([]byte{0xC0, 0x80}), "", ErrInvalidUTF8},
		{"Truncated4Byte", string([]byte{0xF0, 0x9F, 0x92}), "", ErrInvalidUTF8},
		{"InvalidContinuation", string([]byte{0xE2, 0x28, 0xA1}), "", ErrInvalidUTF8},
		{"OutOfRange", string([]byte{0xF4, 0x90, 0x80, 0x80}), "", ErrInvalidUTF8},
		{"TrailingBackslash", "ends with \\", `"ends with \\"`, nil},
		{"TrailingQuote", "ends with \"", `"ends with \""`, nil},
		{"HTMLEscapeLessThan", "<", `"\u003c"`, nil},
		{"HTMLEscapeGreaterThan", ">", `"\u003e"`, nil},
		{"HTMLEscapeAmpersand", "&", `"\u0026"`, nil},
		{"HTMLEscapeAll", "<script>&</script>", `"\u003cscript\u003e\u0026\u003c/script\u003e"`, nil},
		{"HTMLEscapeMixed", "a<b>c&d", `"a\u003cb\u003ec\u0026d"`, nil},
		{"HTMLEscapeWithQuotes", `"<&>"`, `"\"\u003c\u0026\u003e\""`, nil},
		{"U+2029", "\u2029", "\"\\u2029\"", nil},
		{"U+2028", "\u2028", "\"\\u2028\"", nil},
		{"NoOverEscape", "\u2028\u2029", "\"\\u2028\\u2029\"", nil},
	}

	for _, tc := range tests {
		t.Run("append_"+tc.name, func(t *testing.T) {
			got, err := appendString([]byte{}, tc.value)
			Equals(t, tc.wantErr, err)
			Equals(t, tc.want, string(got))
		})
	}
}

func FuzzAppendString(f *testing.F) {
	tests := []string{
		"",
		"simple ascii",
		`quote"slash\test`,
		"control:\b\t\n\f\r",
		"low\u0001\u0002\u001f",
		"こんにちは世界",
		"emoji 😀😅🚀",
		"mixed ascii 日本語 😀 \n \" \\",
		string([]byte{0xEF, 0xBF, 0xBD}),       // U+FFFD replacement character
		string([]byte{0xFF}),                   // invalid UTF-8 byte
		string([]byte{0xE2, 0x82}),             // truncated 3-byte sequence
		string([]byte{0xED, 0xA0, 0x80}),       // surrogate high
		string([]byte{0xED, 0xBF, 0xBF}),       // surrogate low
		string([]byte{0xC0, 0x80}),             // overlong encoding
		string([]byte{0xF0, 0x9F, 0x92}),       // truncated 4-byte
		string([]byte{0xE2, 0x28, 0xA1}),       // invalid continuation
		string([]byte{0xF4, 0x90, 0x80, 0x80}), // out of range
		"ends with \\",
		"ends with \"",
		"<",
		">",
		"&",
		"<script>&</script>",
		"a<b>c&d",
		`"<&>"`,
		"\u2029",
		"\u2028",
		"\u2028\u2029",
	}
	for _, tc := range tests {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		got, err := appendString(nil, s)
		if err != nil {
			if !errors.Is(err, ErrInvalidUTF8) {
				t.Fatal(err)
			}
			return
		}

		want, _ := json.Marshal(s)
		if !bytes.Equal(want, got) {
			t.Fatalf("input=%q\nwant=%q (len=%d)\ngot =%q (len=%d)\nwant hex=%x\ngot hex=%x", s, want, len(want), got, len(got), want, got)
			return
		}
	})
}

func BenchmarkAppendString_Unicode_Throughput(b *testing.B) {
	sizes := stringSize

	rng := rand.New(rand.NewSource(1))
	const batchCount = 1024

	// Prepare samples per size
	samples := make([]string, len(sizes))
	for i, size := range sizes {
		samples[i] = randomString(size, rng)
	}

	buf := make([]byte, 0)

	b.ReportAllocs()
	for i, size := range sizes {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			sample := samples[i]

			totalBytes := size*batchCount*AvgWorstCaseBytesPerString + 2
			buf = make([]byte, 0, totalBytes)

			b.SetBytes(int64(totalBytes))
			b.ResetTimer()

			for b.Loop() {
				buf = buf[:0]
				for j := 0; j < batchCount; j++ {
					buf, _ = appendString(buf, sample)
				}
			}
		})
	}
}

func BenchmarkAppendString_Unicode(b *testing.B) {
	sizes := stringSize

	rng := rand.New(rand.NewSource(1))
	samples := make([]string, len(sizes))
	for i, size := range sizes {
		samples[i] = randomString(size, rng)
	}

	buf := make([]byte, 0)

	b.ReportAllocs()
	for i, size := range sizes {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			sample := samples[i]

			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf = buf[:0]
				buf, _ = appendString(buf, sample)
			}
		})
	}
}

func BenchmarkAppendStringASCII_Throughput(b *testing.B) {
	sizes := stringSize

	rng := rand.New(rand.NewSource(1))
	samples := make([]string, len(sizes))
	for i, size := range sizes {
		samples[i] = randomASCIIString(size, rng)
	}

	buf := make([]byte, 0)
	const batchCount = 1024

	b.ReportAllocs()

	for i, size := range sizes {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			sample := samples[i]

			totalBytes := size*batchCount*AvgWorstCaseBytesPerString + 2
			buf = make([]byte, 0, totalBytes)

			b.SetBytes(int64(totalBytes))
			b.ResetTimer()

			for b.Loop() {
				buf = buf[:0]
				for j := 0; j < batchCount; j++ {
					buf, _ = appendString(buf, sample)
				}
			}
		})
	}
}

func BenchmarkAppendStringASCII(b *testing.B) {
	sizes := stringSize

	rng := rand.New(rand.NewSource(1))
	samples := make([]string, len(sizes))
	for i, size := range sizes {
		samples[i] = randomASCIIString(size, rng)
	}
	var sample string
	buf := make([]byte, 0)

	b.ReportAllocs()
	for i, size := range sizes {
		b.Run("Size_"+strconv.Itoa(size), func(b *testing.B) {
			sample = samples[i]

			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf = buf[:0]
				buf, _ = appendString(buf, sample)
			}
		})
	}
}
