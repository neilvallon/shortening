package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math/bits"
)

type encoder interface {
	Encode(n uint64) []byte
}

type decoder interface {
	Decode(b []byte) (uint64, error)
}

type coder interface {
	encoder
	decoder
}

var (
	InvalidDecodeLen = errors.New("shortening: invalid decode length")
	InvalidCharacter = errors.New("shortening: invalid decode character")
	Overflow         = errors.New("shortening: uint64 overflow")
)

func makeTable(cs string) (t [256]uint8) {
	for i := range t {
		t[i] = 0xFF
	}

	for i, c := range cs {
		t[c] = uint8(i)
	}
	return t
}

func encLen(n uint64, width int, table []uint64) int {
	bit := bits.Len64(n) / width
	if table[bit] <= n {
		bit++
	}

	return bit
}
