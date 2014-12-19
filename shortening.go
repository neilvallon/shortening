package shortening

import (
	"bytes"
	"errors"
	"math"
)

var (
	charSet = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	pow64   = []uint64{
		1,
		64,
		64 * 64,
		64 * 64 * 64,
		64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64 * 64 * 64 * 64 * 64,
		64 * 64 * 64 * 64 * 64 * 64 * 64 * 64 * 64 * 64,
	}
)

// Encode turns an uint64 into a slice of characters from 'charSet'
func Encode(n uint64) []byte {
	if n == math.MaxUint64 {
		return []byte("_---------O")
	}

	b := make([]byte, 0, 11) // preallocate to avoid growslice

	n++
	for 0 != n {
		n--
		b = append(b, charSet[n%64])
		n /= 64
	}

	return b
}

// Decode turns a slice of characters back into the original unit64.
//
// Errors are returned for invalid characters or input that would
// cause an overflow.
func Decode(b []byte) (n uint64, err error) {
	if 11 < len(b) {
		return 0, errors.New("shortening: too many bytes to decode")
	}

	for i, c := range b {
		ind := bytes.IndexByte(charSet, c) // Faster than using a map
		if ind == -1 {
			return 0, errors.New("shortening: invalid decode character")
		}

		nn := n + uint64(ind+1)*pow64[i]
		if nn-1 < n {
			return 0, errors.New("shortening: int64 overflow")
		}

		n = nn
	}
	return n - 1, nil
}
