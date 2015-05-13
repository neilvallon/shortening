package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math"
)

var (
	// charSet - 64 ascii byte values (0-127)
	charSet     = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	lookupTable = makeTable(charSet)
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
	if 11 < len(b) || len(b) == 0 {
		return 0, errors.New("shortening: invalid decode length")
	}

	for i, c := range b {
		ind := lookupTable[c]
		if ind == -1 {
			return 0, errors.New("shortening: invalid decode character")
		}

		nn := n + uint64(ind+1)<<uint(6*i)
		if nn-1 < n {
			return 0, errors.New("shortening: int64 overflow")
		}

		n = nn
	}

	return n - 1, nil
}

func makeTable(cs []byte) []int8 {
	t := make([]int8, 256)

	// fill table with error values
	for i := range t {
		t[i] = -1
	}

	for i, c := range cs {
		t[c] = int8(i)
	}

	return t
}
