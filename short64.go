package shortening // import "vallon.me/shortening"

import (
	"errors"
)

var (
	charSet     = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	lookupTable = makeTable(charSet)

	lookupBigTable = makeBigTable(charSet)
)

const (
	// Min values are the smallest value of 'n' that uses the number of
	// place values.
	min02 = 1 << 6
	min03 = (min02 + 1) << 6
	min04 = (min03 + 1) << 6
	min05 = (min04 + 1) << 6
	min06 = (min05 + 1) << 6
	min07 = (min06 + 1) << 6
	min08 = (min07 + 1) << 6
	min09 = (min08 + 1) << 6
	min10 = (min09 + 1) << 6
	min11 = (min10 + 1) << 6
)

var minTable = [...]uint64{0, min02, min03, min04, min05, min06, min07, min08, min09, min10, min11}

// Encode turns an uint64 into a slice of characters from 'charSet'
func Encode(n uint64) []byte {
	var buf [11]byte

	nn := n - min11
	for i, m := range minTable {
		if n < m {
			return buf[:i]
		}

		buf[i], nn = charSet[nn&63], nn>>6
	}

	return buf[:]
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
		n = n<<6 | uint64(ind)
	}

	n += minTable[len(b)-1]

	if invalid == 0xFF {
		return 0, errors.New("shortening: invalid decode character")
	}

	// 1171221845949812800 is the minimum value to have len == 11
	// any lower value is an overflow.
	if len(b) == 11 && n < min11 {
		return 0, errors.New("shortening: int64 overflow")
	}

	return n, nil
}

func makeTable(cs []byte) (t [256]uint8) {
	for i := range t {
		t[i] = 0xFF
	}

	for i, c := range cs {
		t[c] = uint8(i)
	}
	return t
}
