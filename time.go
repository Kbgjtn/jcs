package jcs

import (
	"time"
)

// appendTime appends a time.Time value to dst as a JSON string.
//
// The time is first converted to UTC and formatted using RFC3339Nano,
// producing a deterministic, canonical representation such as
// "2019-01-28T07:45:10Z". Per RFC 8785 (JSON Canonicalization Scheme),
// time values are treated as ordinary JSON strings with no special
// normalization beyond consistent formatting. The resulting string is
// then quoted and escaped using appendString to ensure valid JSON.

// appendTime appends an RFC3339 or RFC3339Nano timestamp in UTC to dst and returns dst.
func appendTime(dst []byte, t time.Time) []byte {
	// Use UTC for canonical JSON timestamps
	dst = append(dst, '"')
	tt := t.UTC()

	year, month, day := tt.Date()
	hour, min, sec := tt.Clock()
	nsec := tt.Nanosecond()

	// Reserve capacity: 20 for RFC3339, up to 30 for RFC3339Nano
	need := 20
	if nsec != 0 {
		// RFC3339Nano can be up to 29 chars like 2006-01-02T15:04:05.123456789Z
		need = 20 + 1 + 9 // dot + up to 9 digits
	}
	if cap(dst)-len(dst) < need {
		newCap := cap(dst)*2 + need
		if newCap < len(dst)+need {
			newCap = len(dst) + need
		}
		new := make([]byte, len(dst), newCap)
		copy(new, dst)
		dst = new
	}

	// Append date and time: YYYY-MM-DDTHH:MM:SS
	dst = appendOffsetUTCZeroFourDigit(dst, year)
	dst = append(dst, '-')
	dst = appendOffsetUTCTwoDigit(dst, int(month))
	dst = append(dst, '-')
	dst = appendOffsetUTCTwoDigit(dst, day)
	dst = append(dst, 'T')
	dst = appendOffsetUTCTwoDigit(dst, hour)
	dst = append(dst, ':')
	dst = appendOffsetUTCTwoDigit(dst, min)
	dst = append(dst, ':')
	dst = appendOffsetUTCTwoDigit(dst, sec)

	if nsec != 0 {
		// Trim trailing zeros in fractional seconds to match RFC3339Nano semantics
		// produce exactly the digits needed (no trailing zeros)
		// write '.' then up to 9 digits
		frac := nsec
		// produce 9 digits into a small array then trim trailing zeros
		var buf [9]byte
		for i := 8; i >= 0; i-- {
			buf[i] = byte('0' + frac%10)
			frac /= 10
		}
		// find last non-zero
		last := 8
		for last >= 0 && buf[last] == '0' {
			last--
		}
		if last >= 0 {
			dst = append(dst, '.')
			dst = append(dst, buf[:last+1]...)
		}
	}

	dst = append(dst, 'Z')
	dst = append(dst, '"')
	return dst
}

// appendOffsetUTCTwoDigit appends a zero-padded 2-digit decimal of v (0..99).
func appendOffsetUTCTwoDigit(dst []byte, v int) []byte {
	dst = append(dst, byte('0'+(v/10)), byte('0'+(v%10)))
	return dst
}

// appendOffsetUTCZeroFourDigit appends a zero-padded 4-digit decimal of v (year).
func appendOffsetUTCZeroFourDigit(dst []byte, v int) []byte {
	// v is typically >= 0 and small (e.g., 2026)
	d1 := v / 1000
	d2 := (v / 100) % 10
	d3 := (v / 10) % 10
	d4 := v % 10
	dst = append(dst, byte('0'+d1), byte('0'+d2), byte('0'+d3), byte('0'+d4))
	return dst
}
