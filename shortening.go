package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math"
)

var (
	charSet     = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	pow64       = powerArray(64)
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
		if len(lookupTable) <= int(c) {
			return 0, errors.New("shortening: invalid decode character")
		}

		ind := lookupTable[c]
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

func makeTable(cs []byte) []int {
	t := make([]int, max(cs)+1)

	// fill table with error values
	for i := range t {
		t[i] = -1
	}

	for i, c := range cs {
		t[c] = i
	}

	return t
}

func max(b []byte) byte {
	var n byte = 0x00

	for _, c := range b {
		if c > n {
			n = c
		}
	}

	return n
}

func powerArray(base uint64) []uint64 {
	parr := make([]uint64, 1, 11)
	parr[0] = 1

	for i := 1; ; i++ {
		n := base * parr[i-1]

		if n < parr[i-1] {
			break // overflow
		}

		parr = append(parr, n)
	}

	return parr
}
