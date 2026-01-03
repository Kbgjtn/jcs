package jcs

import (
	"slices"
	"sort"
)

// appendObject serializes a map[string]any (JSON object) into the destination byte slice `dst`.
// The function sorts the keys lexicographically, processes UTF-16 encoding for key/value pairs,
// and ensures the JSON object is serialized in canonical form as per RFC 8785.
func appendObject(dst []byte, obj map[string]any) ([]byte, error) {
	dstLen := len(dst)
	dst = append(dst, '{')
	if len(obj) == 0 {
		return append(dst, '}'), nil
	}

	totalBytesKey := 0
	totalBytesValues := 0
	for k, v := range obj {
		totalBytesKey += utf16Len(k)
		totalBytesValues += estimateValueSize(v)
	}

	// shared UTF-16 buffer for all keys
	// heuristic: ASCII keys dominate - ~1 code unit per byte
	utf16buf := make([]uint16, 0, totalBytesKey)

	keys := make([]kv, len(obj))
	i := 0
	for k := range obj {
		start := len(utf16buf)
		var n int
		var err error

		utf16buf, n, err = appendUTF16(utf16buf, k)
		if err != nil {
			return dst[:dstLen], err
		}

		keys[i] = kv{
			raw:   k,
			len:   uint32(n),
			start: uint32(start),
		}
		i++
	}

	sort.Sort(kvSlice{keys: keys, utf16buf: utf16buf})

	estimatedTotalBytes := totalBytesKey*2 + totalBytesValues + len(obj)*2
	dst = slices.Grow(dst, estimatedTotalBytes)

	for i, key := range keys {
		if i > 0 {
			dst = append(dst, ',')
		}

		var err error

		// key

		// dst = slices.Grow(dst, len(key.raw)*6+2)
		dst, err = appendString(dst, key.raw)
		if err != nil {
			return dst[:dstLen], err
		}

		dst = append(dst, ':')
		dst, err = Append(dst, obj[key.raw])
		if err != nil {
			return dst[:dstLen], err
		}
	}

	dst = append(dst, '}')
	return dst, nil
}

func estimateValueSize(v interface{}) int {
	switch val := v.(type) {
	case string:
		// Worst case: every character could be escaped
		return len(val)*2 + 2 // quotes + escape
	case bool:
		return 5 // max for 'false'
	case nil:
		return 4
	case int, int32, int64, float32, float64:
		return 20 // max digits for large numbers
	case []interface{}:
		size := 2 // [ ]
		for _, e := range val {
			size += estimateValueSize(e) + 1 // comma
		}
		if len(val) > 0 {
			size-- // remove last comma
		}
		return size
	case map[string]interface{}:
		size := 2 // { }
		for k, e := range val {
			size += len(k)*2 + 2 + 1         // quotes + colon
			size += estimateValueSize(e) + 1 // comma
		}
		if len(val) > 0 {
			size-- // remove last comma
		}
		return size
	default:
		return 32 // fallback for unknown types
	}
}
