package jcs

import (
	"unicode/utf16"
	"unicode/utf8"
)

// kv represents a key stored as a UTF-16 encoded string segment.
// It stores the raw string, along with the starting position and length of the key
// in terms of UTF-16 code units (not bytes), which is important for accurate string
// processing and comparisons involving UTF-16 encoded data.
type kv struct {
	// raw holds the original UTF-16 encoded string.
	// It's not necessarily a UTF-16 encoded byte sequence, but a regular Go string
	// interpreted or processed as UTF-16 for the purpose of key management.
	raw string

	// start is the starting position of the key within some larger string or dataset.
	// It helps to know where the key begins when working with a collection of UTF-16 data.
	start uint32

	// len is the length of the key in terms of UTF-16 code units (i.e., number of 16-bit units).
	// This accounts for surrogate pairs and ensures proper indexing.
	len uint32
}

type kvSlice struct {
	keys     []kv
	utf16buf []uint16
}

func (kvs kvSlice) Len() int {
	return len(kvs.keys)
}

func (kvs kvSlice) Less(i, j int) bool {
	aStart, aLen := int(kvs.keys[i].start), kvs.keys[i].len
	bStart, bLen := int(kvs.keys[j].start), kvs.keys[j].len

	n := int(min(aLen, bLen))
	for x := 0; x < n; x++ {
		if kvs.utf16buf[aStart+x] != kvs.utf16buf[bStart+x] {
			return kvs.utf16buf[aStart+x] < kvs.utf16buf[bStart+x]
		}
	}
	return aLen < bLen
}

func (s kvSlice) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
}

// appendUTF16 converts a UTF-8 encoded string into UTF-16 code units.
//
// The function iterates over the input string using utf8.DecodeRuneInString,
// decoding one rune at a time and appending its UTF-16 representation to buf.
// Runes in the Basic Multilingual Plane (U+0000–U+FFFF) are encoded as a single
// 16-bit value. Supplementary characters (U+10000 and above) are encoded as a
// surrogate pair using utf16.EncodeRune.
//
// Invalid UTF-8 handling:
//   - utf8.DecodeRuneInString returns utf8.RuneError (U+FFFD) when it encounters
//     malformed input.
//   - If RuneError is returned with size == 1, this indicates an invalid single
//     byte sequence (such as 0xFF). In this case, appendUTF16 returns ErrInvalidUTF8.
//   - Other cases of RuneError (size > 1) are treated as replacement characters
//     and encoded as U+FFFD. This matches Go’s convention of replacing malformed
//     multi-byte sequences with U+FFFD rather than failing.
//
// The returned slice contains the appended UTF-16 code units, the number of
// units written, and an error if invalid UTF-8 was detected.
func appendUTF16(buf []uint16, s string) ([]uint16, int, error) {
	start := len(buf)

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			return buf, 0, ErrInvalidUTF8
		}

		// handle the encoding into UTF-16
		// check if the rune is less than Basic Multilingual Plane (BMP, [U+0000–U+FFFF])
		if r < 0x10_000 {
			// means r can be represented with a single 16-bit value in UTF-16
			buf = append(buf, uint16(r))
		} else {
			// Since UTF-16 uses 16-bit units, it cannot directly represent these
			// higher code points with a single unit. Instead, it uses two 16-bit
			// code units (a high and lower surrogate) to encode one character.
			// basically, we need a way to fit them into two 16-bit units.
			//
			// Surrogate range:
			//  - High surrogates: U+D800–U+DBFF
			//  - Low surrogates: U+DC00–U+DFFF
			// these ranges are reserved exclusively for surrogate use and are not
			// valid standalone characters.

			// substracting into 20 bits padded
			// and split into two 10‑bit halves
			// we get the result of high 10-bits and low 10-bits
			h, l := utf16.EncodeRune(r)
			buf = append(buf, uint16(h), uint16(l))
		}
		i += size
	}

	return buf, len(buf) - start, nil
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		if r <= 0xFFFF {
			n++
		} else {
			n += 2
		}
	}
	return n
}
