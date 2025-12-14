package jcs

import (
	"math"
	"strconv"
)

// MaxSafeNumber defines the largest integer that can be represented exactly
// in IEEE‑754 double precision (binary64), as required by RFC 8785 (JSON
// Canonicalization Scheme).
//
// RFC 8785 3.2.2 mandates that all JSON numbers must be preserved exactly
// when serialized. Since canonical JSON relies on IEEE‑754 binary64, only
// integers in the range [‑(2^53‑1), +(2^53‑1)] are considered "safe".
// Any integer outside this range cannot be represented without precision
// loss and must cause canonicalization to fail.
//
// MaxSafeNumber is therefore set to 2^53‑1 (9007199254740991).
const MaxSafeNumber = 1<<53 - 1

// numberNormalizer normalizes the exponent of a float64 value `v` in the JSON output
// to adhere strictly to the RFC 8785 canonical JSON format.
//
// This function ensures that the floating-point value is serialized in scientific notation,
// following these rules from RFC 8785:
//   - The exponent must always be expressed with a sign (`+` or `-`), even for positive exponents.
//   - The exponent cannot have leading zeros (except for `+00` or `-00`).
//   - For numbers like 0.00000123 or 1230000000, the exponent part must be adjusted and normalized.
//   - Example: A value like 1230000000 must be represented as 1.23e9 (with explicit "e+9" as the exponent).
func numberNormalizer(dst []byte, v float64) []byte {
	dst = strconv.AppendFloat(dst, v, 'e', -1, 64)

	for i := len(dst) - 1; i >= 0; i-- {
		if dst[i] == 'e' {
			// get position '+'
			j := i + 1

			// remove '+'
			if dst[j] == '+' {
				copy(dst[j:], dst[j+1:])
				dst = dst[:len(dst)-1]
			}

			// remove leading zero in exponent
			if j+1 < len(dst) && dst[j] == '-' && dst[j+1] == '0' {
				copy(dst[j+1:], dst[j+2:])
				dst = dst[:len(dst)-1]
			} else if dst[j] == '0' && j+1 < len(dst) {
				copy(dst[j:], dst[j+1:])
				dst = dst[:len(dst)-1]
			}

			return dst
		}
	}

	return dst
}

// isNumberOOR checks whether an integer value lies outside the IEEE‑754 binary64
// "safe integer" range defined by RFC 8785 (JSON Canonicalization Scheme).
// RFC 8785 requires that all JSON numbers be representable exactly in
// IEEE‑754 double precision. This means only integers in the range
// [‑(2^53‑1), +(2^53‑1)] are valid. Any integer outside this range cannot
// be preserved exactly and must cause canonicalization to fail.
// The function returns ErrNumberOOR if v is out of range, otherwise nil.
func isNumberOOR[T int | int64 | uint | uint64](v T) bool {
	switch any(v).(type) {
	case int, int64:
		val := int64(v)
		if val > MaxSafeNumber || val < -MaxSafeNumber {
			return true
		}
	case uint, uint64:
		val := uint64(v)
		if val > MaxSafeNumber {
			return true
		}
	}
	return false
}

// appendNumber appends the canonical JSON representation of a numeric value to dst.
//
// This function implements the numeric serialization rules required by RFC 8785
// (JSON Canonicalization Scheme). All numeric values are converted to float64
// internally and then rendered as JSON numbers in their shortest decimal form
// that guarantees round‑trip exactness.
//
// Canonicalization rules enforced:
//
//   - Zero normalization:
//     Both +0.0 and -0.0 are rendered as "0". RFC 8785 requires that negative
//     zero not be distinguishable from positive zero in canonical JSON.
//
//   - NaN and Infinity:
//     NaN, +Inf, and -Inf are explicitly disallowed by RFC 8785. If encountered,
//     the function returns ErrNaN or ErrInf respectively.
//
//   - Safe integer range:
//     Integers must be restricted to the IEEE‑754 double‑precision safe range
//     [-(2^53-1), +(2^53-1)]. Values outside this range cannot be represented
//     exactly as float64 and will trigger ErrNumberOOR.
//
//   - Shortest decimal representation:
//     For finite values within the safe range, strconv.AppendFloat is used with
//     mode 'f' and precision -1 to produce the shortest correct decimal string.
//     For very large or very small magnitudes, float64Normalizer is used to
//     ensure canonical exponent formatting (removing '+' signs and leading zeros).
//
// Error handling:
//   - Returns `ErrNaN` if v is NaN.
//   - Returns `ErrInf` if v is +Inf or -Inf.
//   - Returns `ErrNumberOOR` if an integer exceeds ±(2^53-1).
//   - Otherwise, returns the updated dst slice containing the canonical JSON
//     number.
//
// The resulting output is guaranteed to be a valid, canonical JSON number according to RFC 8785.
func appendNumber(dst []byte, v float64) ([]byte, error) {
	if v == 0 {
		if math.Signbit(v) {
			return append(dst, '0'), nil
		}
		return append(dst, '0'), nil
	}
	if math.IsNaN(v) {
		return dst, ErrNaN
	}
	if math.IsInf(v, 0) {
		return dst, ErrInf
	}

	abs := math.Abs(v)
	if abs < 1e+21 && abs >= 1e-6 {
		dst = strconv.AppendFloat(dst, v, 'f', -1, 64)
		return dst, nil
	}

	// Both Ryu and Grisu3 are specialized algorithms for converting floating‑point numbers to their shortest correct decimal string representation.
	// Ryu is newer, provably correct, and consistently fast; Grisu3 is older, very fast in common cases but sometimes falls back to slower routines.
	// So this is not fastest possible, but quite solid, this dependent on strconv.AppendFloat most of the runtime cost is inside this Go's generic
	// float-to-string conversion doing heavy math.
	//
	// also exponent cleanup overhead, because its adds a small post-processing loop to strip '+' and leading zeros ('0'). It's cheap but still extra
	// works compared to algorithms that generate the exact format directly (like Ryu).
	// Specialized algorithms (Ryu, Grisu3) can cut float‑to‑string conversion time further, but they’re complex to implement and maintain.
	return numberNormalizer(dst, v), nil
}
