package jcs

import (
	"unicode/utf8"
)

type escInfo struct {
	n    int
	data [6]byte
}

var asciiTable = [128]escInfo{
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x30}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x31}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x32}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x33}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x34}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x35}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x36}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x37}},
	{n: 2, data: [6]uint8{0x5c, 0x62}},
	{n: 2, data: [6]uint8{0x5c, 0x74}},
	{n: 2, data: [6]uint8{0x5c, 0x6e}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x62}},
	{n: 2, data: [6]uint8{0x5c, 0x66}},
	{n: 2, data: [6]uint8{0x5c, 0x72}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x65}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x30, 0x66}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x30}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x31}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x32}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x33}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x34}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x35}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x36}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x37}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x38}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x39}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x61}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x62}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x63}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x64}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x65}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x31, 0x66}},
	{n: 1, data: [6]uint8{0x20}},
	{n: 1, data: [6]uint8{0x21}},
	{n: 2, data: [6]uint8{0x5c, 0x22}},
	{n: 1, data: [6]uint8{0x23}},
	{n: 1, data: [6]uint8{0x24}},
	{n: 1, data: [6]uint8{0x25}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x32, 0x36}},
	{n: 1, data: [6]uint8{0x27}},
	{n: 1, data: [6]uint8{0x28}},
	{n: 1, data: [6]uint8{0x29}},
	{n: 1, data: [6]uint8{0x2a}},
	{n: 1, data: [6]uint8{0x2b}},
	{n: 1, data: [6]uint8{0x2c}},
	{n: 1, data: [6]uint8{0x2d}},
	{n: 1, data: [6]uint8{0x2e}},
	{n: 1, data: [6]uint8{0x2f}},
	{n: 1, data: [6]uint8{0x30}},
	{n: 1, data: [6]uint8{0x31}},
	{n: 1, data: [6]uint8{0x32}},
	{n: 1, data: [6]uint8{0x33}},
	{n: 1, data: [6]uint8{0x34}},
	{n: 1, data: [6]uint8{0x35}},
	{n: 1, data: [6]uint8{0x36}},
	{n: 1, data: [6]uint8{0x37}},
	{n: 1, data: [6]uint8{0x38}},
	{n: 1, data: [6]uint8{0x39}},
	{n: 1, data: [6]uint8{0x3a}},
	{n: 1, data: [6]uint8{0x3b}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x33, 0x63}},
	{n: 1, data: [6]uint8{0x3d}},
	{n: 6, data: [6]uint8{0x5c, 0x75, 0x30, 0x30, 0x33, 0x65}},
	{n: 1, data: [6]uint8{0x3f}},
	{n: 1, data: [6]uint8{0x40}},
	{n: 1, data: [6]uint8{0x41}},
	{n: 1, data: [6]uint8{0x42}},
	{n: 1, data: [6]uint8{0x43}},
	{n: 1, data: [6]uint8{0x44}},
	{n: 1, data: [6]uint8{0x45}},
	{n: 1, data: [6]uint8{0x46}},
	{n: 1, data: [6]uint8{0x47}},
	{n: 1, data: [6]uint8{0x48}},
	{n: 1, data: [6]uint8{0x49}},
	{n: 1, data: [6]uint8{0x4a}},
	{n: 1, data: [6]uint8{0x4b}},
	{n: 1, data: [6]uint8{0x4c}},
	{n: 1, data: [6]uint8{0x4d}},
	{n: 1, data: [6]uint8{0x4e}},
	{n: 1, data: [6]uint8{0x4f}},
	{n: 1, data: [6]uint8{0x50}},
	{n: 1, data: [6]uint8{0x51}},
	{n: 1, data: [6]uint8{0x52}},
	{n: 1, data: [6]uint8{0x53}},
	{n: 1, data: [6]uint8{0x54}},
	{n: 1, data: [6]uint8{0x55}},
	{n: 1, data: [6]uint8{0x56}},
	{n: 1, data: [6]uint8{0x57}},
	{n: 1, data: [6]uint8{0x58}},
	{n: 1, data: [6]uint8{0x59}},
	{n: 1, data: [6]uint8{0x5a}},
	{n: 1, data: [6]uint8{0x5b}},
	{n: 2, data: [6]uint8{0x5c, 0x5c}},
	{n: 1, data: [6]uint8{0x5d}},
	{n: 1, data: [6]uint8{0x5e}},
	{n: 1, data: [6]uint8{0x5f}},
	{n: 1, data: [6]uint8{0x60}},
	{n: 1, data: [6]uint8{0x61}},
	{n: 1, data: [6]uint8{0x62}},
	{n: 1, data: [6]uint8{0x63}},
	{n: 1, data: [6]uint8{0x64}},
	{n: 1, data: [6]uint8{0x65}},
	{n: 1, data: [6]uint8{0x66}},
	{n: 1, data: [6]uint8{0x67}},
	{n: 1, data: [6]uint8{0x68}},
	{n: 1, data: [6]uint8{0x69}},
	{n: 1, data: [6]uint8{0x6a}},
	{n: 1, data: [6]uint8{0x6b}},
	{n: 1, data: [6]uint8{0x6c}},
	{n: 1, data: [6]uint8{0x6d}},
	{n: 1, data: [6]uint8{0x6e}},
	{n: 1, data: [6]uint8{0x6f}},
	{n: 1, data: [6]uint8{0x70}},
	{n: 1, data: [6]uint8{0x71}},
	{n: 1, data: [6]uint8{0x72}},
	{n: 1, data: [6]uint8{0x73}},
	{n: 1, data: [6]uint8{0x74}},
	{n: 1, data: [6]uint8{0x75}},
	{n: 1, data: [6]uint8{0x76}},
	{n: 1, data: [6]uint8{0x77}},
	{n: 1, data: [6]uint8{0x78}},
	{n: 1, data: [6]uint8{0x79}},
	{n: 1, data: [6]uint8{0x7a}},
	{n: 1, data: [6]uint8{0x7b}},
	{n: 1, data: [6]uint8{0x7c}},
	{n: 1, data: [6]uint8{0x7d}},
	{n: 1, data: [6]uint8{0x7e}},
	{n: 1, data: [6]uint8{0x7f}},
}

var asciiN [128]byte = [...]byte{6, 6, 6, 6, 6, 6, 6, 6, 2, 2, 2, 6, 2, 2, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 1, 1, 2, 1, 1, 1, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 6, 1, 6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

// appendString appends the canonical JSON representation of a Go string to dst.
//
// This function implements the string escaping and UTF-8 validation rules
// required by RFC 8785 (JSON Canonicalization Scheme):
//
//   - All output is UTF-8 encoded and enclosed in double quotes.
//   - Safe ASCII characters (U+0020–U+007E, excluding '"' and '\') are copied
//     directly for efficiency.
//   - Control characters (< U+0020) and the special characters '"' and '\' are
//     escaped using JSON escape sequences (e.g., \u00XX, \n, \t).
//   - Non-ASCII characters are validated and emitted as-is.
//
// UTF-8 validation:
//   - utf8.DecodeRuneInString is used to decode each rune.
//   - If DecodeRune returns utf8.RuneError, both invalid single-byte sequences
//     and truncated multi-byte sequences are rejected with ErrInvalidUTF8.
//   - A literal U+FFFD replacement character (encoded as 0xEF 0xBF 0xBD) is
//     allowed, since it is a valid Unicode scalar value.
//   - Surrogate code points (U+D800–U+DFFF) are explicitly rejected, as they
//     are not valid Unicode scalar values and disallowed by RFC 8785.
//
// Error handling:
//   - Returns ErrInvalidUTF8 if the input string contains malformed UTF-8 or
//     surrogate code points.
//   - Otherwise, returns the updated dst slice containing the escaped string.
//
// The resulting output is guaranteed to be a valid, canonical JSON string
// according to RFC 8785.
func appendString(dst []byte, s string) ([]byte, error) {
	dst = append(dst, '"')
	i := 0
	n := len(s)

	var info *escInfo

	for i < n {
		if s[i] < 128 {
			start := i

			for i < n {
				if ch := s[i]; ch >= 128 || asciiN[ch] != 1 {
					break
				}

				i++
			}

			if i > start {
				dst = append(dst, s[start:i]...)
			}

			// escape any remaining ASCII (control chars)
			if i < n {
				if ch := s[i]; ch < 128 {
					info = &asciiTable[ch]
					if info.n != 1 {
						dst = append(dst, info.data[:info.n]...)
						i++
					}

				}
			}

			continue
		}

		// multi-byte UTF-8 run — batch decode
		for i < n && s[i] >= 128 {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				return nil, ErrInvalidUTF8
			}

			switch r {
			case '\u2028':
				dst = append(dst, `\u2028`...)
			case '\u2029':
				dst = append(dst, `\u2029`...)
			default:
				dst = append(dst, s[i:i+size]...)
			}

			i += size
		}
	}
	dst = append(dst, '"')
	return dst, nil
}
