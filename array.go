package jcs

// appendSlice appends the canonical JSON representation of a Go slice to dst.
//
// This function implements array serialization rules required by RFC 8785
// (JSON Canonicalization Scheme). It ensures that slices are encoded as JSON
// arrays with elements serialized in canonical form.
//
// Canonicalization rules enforced:
//
//   - Arrays are always enclosed in square brackets '[' and ']'.
//   - Elements are serialized in their natural order without reordering.
//   - Elements are separated by a single comma ',' with no extra whitespace.
//   - Each element is encoded using Append, which applies the appropriate
//     canonicalization rules for its type (string, number, boolean, object, etc.).
//   - If any element encoding fails, the function restores dst to its original
//     length and returns the error.
//
// Error handling:
//   - Returns any error produced by Append when encoding an element.
//   - Common errors include ErrUnsupportedType, ErrNaN, ErrInf, ErrInvalidUTF8,
//     or ErrNumberOOR depending on the element type.
//   - If an error occurs, no partial array is returned; dst is reset to its
//     original state before appendSlice was called.
//
// The resulting output is guaranteed to be a valid, canonical JSON array
// according to RFC 8785, with each element individually validated and encoded.
func appendSlice[T any](dst []byte, arr []T) ([]byte, error) {
	dstLen := len(dst)
	dst = append(dst, '[')

	for i, v := range arr {
		if i > 0 {
			dst = append(dst, ',')
		}

		var err error
		dst, err = Append(dst, v)
		if err != nil {
			dst = dst[:dstLen]
			return dst, err
		}
	}

	dst = append(dst, ']')
	return dst, nil
}
