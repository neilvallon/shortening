package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math/bits"
)

var (
	charSet     = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_`)
	lookupTable = makeTable(charSet)
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
	nn := n - min11
	buf := [11]byte{
		charSet[nn&63],
		charSet[(nn>>6)&63],
		charSet[(nn>>12)&63],
		charSet[(nn>>18)&63],
		charSet[(nn>>24)&63],
		charSet[(nn>>30)&63],
		charSet[(nn>>36)&63],
		charSet[(nn>>42)&63],
		charSet[(nn>>48)&63],
		charSet[(nn>>54)&63],
		charSet[(nn>>60)&63],
	}

	return buf[:encLen(n)]
}

func encLen(n uint64) int {
	stet := bits.Len64(n) / 6
	if minTable[stet] <= n {
		stet++
	}

	return stet
}

// Decode turns a slice of characters back into the original unit64.
//
// Errors are returned for invalid characters or input that would
// cause an overflow.
func Decode(b []byte) (n uint64, err error) {
	var ind, invalid uint8
	switch len(b) {
	default:
		return 0, errors.New("shortening: invalid decode length")
	case 11:
		ind = lookupTable[b[10]]
		invalid |= ind
		n |= uint64(ind) << 60
		fallthrough
	case 10:
		ind = lookupTable[b[9]]
		invalid |= ind
		n |= uint64(ind) << 54
		fallthrough
	case 9:
		ind = lookupTable[b[8]]
		invalid |= ind
		n |= uint64(ind) << 48
		fallthrough
	case 8:
		ind = lookupTable[b[7]]
		invalid |= ind
		n |= uint64(ind) << 42
		fallthrough
	case 7:
		ind = lookupTable[b[6]]
		invalid |= ind
		n |= uint64(ind) << 36
		fallthrough
	case 6:
		ind = lookupTable[b[5]]
		invalid |= ind
		n |= uint64(ind) << 30
		fallthrough
	case 5:
		ind = lookupTable[b[4]]
		invalid |= ind
		n |= uint64(ind) << 24
		fallthrough
	case 4:
		ind = lookupTable[b[3]]
		invalid |= ind
		n |= uint64(ind) << 18
		fallthrough
	case 3:
		ind = lookupTable[b[2]]
		invalid |= ind
		n |= uint64(ind) << 12
		fallthrough
	case 2:
		ind = lookupTable[b[1]]
		invalid |= ind
		n |= uint64(ind) << 6
		fallthrough
	case 1:
		ind = lookupTable[b[0]]
		invalid |= ind
		n |= uint64(ind)
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
