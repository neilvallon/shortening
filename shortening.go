package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math"
)

var (
	charSet     = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	lookupTable = makeTable(charSet)
)

// Encode turns an uint64 into a slice of characters from 'charSet'
func Encode(n uint64) []byte {
	var buf [11]byte
	var i int

	for {
		buf[i], n = charSet[n&63], (n>>6)-1
		i++

		if n == math.MaxUint64 {
			return buf[:i]
		}
	}
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
		if ind == 0 {
			return 0, errors.New("shortening: invalid decode character")
		}

		nn := n + uint64(ind)<<uint(6*i)
		if nn-1 < n {
			return 0, errors.New("shortening: int64 overflow")
		}

		n = nn
	}

	return n - 1, nil
}

func makeTable(cs []byte) (t [256]uint8) {
	for i, c := range cs {
		t[c] = uint8(i + 1)
	}
	return t
}
