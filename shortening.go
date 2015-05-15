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

	var invalid uint8
	for i := len(b) - 1; 0 <= i; i-- {
		ind := lookupTable[b[i]]
		invalid |= ind
		n = n<<6 + uint64(ind)
	}

	if invalid == 0xFF {
		return 0, errors.New("shortening: invalid decode character")
	}

	// 1171221845949812800 is the minimum value to have len == 11
	// any lower value is an overflow.
	if len(b) == 11 && n-1 < 1171221845949812800 {
		return 0, errors.New("shortening: int64 overflow")
	}

	return n - 1, nil
}

func makeTable(cs []byte) (t [256]uint8) {
	for i := range t {
		t[i] = 0xFF
	}

	for i, c := range cs {
		t[c] = uint8(i + 1)
	}
	return t
}
