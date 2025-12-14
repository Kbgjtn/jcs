package jcs

import (
	"bytes"
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
func appendTime(dst []byte, t time.Time) []byte {
	// normalize to UTC
	t = t.UTC()

	// Format with nanosecond
	dst = append(dst, '"')
	dst = t.AppendFormat(dst, time.RFC3339Nano)

	// Trim trailing zeros in fractional seconds if present
	if i := bytes.IndexByte(dst, '.'); i != -1 {
		// find end of fractional part before 'Z' or '+'
		end := len(dst) - 1
		for j := end - 1; j > i; j-- {
			if dst[j] == '0' {
				end--
			} else {
				break
			}
		}
		// if we trimmed all fractional digits, also trim the '.'
		if dst[end-1] == '.' {
			end--
		}
		dst = dst[:end+1] // keep suffix (Z or +00:00)
	}

	// Replace "+00:00" with "Z" for UTC
	if bytes.HasSuffix(dst, []byte("+00:00")) {
		dst = append(dst[:len(dst)-6], 'Z')
	}

	dst = append(dst, '"')

	return dst
}
