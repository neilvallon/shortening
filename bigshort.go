package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math/big"
)

var (
	big1  = big.NewInt(1)
	big64 = big.NewInt(64)
)

// EncodeBig return a text representation of any positive integer
// using a base64 character set.
func EncodeBig(n *big.Int) ([]byte, error) {
	if n.Sign() == -1 {
		return nil, errors.New("shortening: encode values must be positive")
	}

	buf := make([]byte, 0, 8)
	var i int

	for {
		var m big.Int
		n.DivMod(n, big64, &m)

		buf = append(buf, charSet[m.Int64()])
		i++

		if n.Sign() == 0 {
			return buf, nil
		}

		n.Sub(n, big1)
	}
}

// DecodeBig reverses an encoded string into a number.
func DecodeBig(b []byte) (*big.Int, error) {
	if len(b) == 0 {
		return nil, errors.New("shortening: invalid decode length")
	}

	var n big.Int
	for i := len(b) - 1; 0 <= i; i-- {
		ind := lookupBigTable[b[i]]
		if ind == nil {
			return nil, errors.New("shortening: invalid decode character")
		}

		n.Mul(&n, big64)
		n.Add(&n, ind)
	}

	return n.Sub(&n, big1), nil
}

func makeBigTable(cs []byte) (t [256]*big.Int) {
	for i, c := range cs {
		t[c] = big.NewInt(int64(i + 1))
	}
	return t
}
