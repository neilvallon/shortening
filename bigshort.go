package shortening // import "vallon.me/shortening"

import (
	"errors"
	"math/big"
)

var (
	lookupBigTable = makeBigTable(CharSet64)

	big1  = big.NewInt(1)
	big63 = big.NewInt(63)
)

// EncodeBig return a text representation of any positive integer
// using a base64 character set.
func EncodeBig(np *big.Int) ([]byte, error) {
	if np.Sign() == -1 {
		return nil, errors.New("shortening: encode values must be positive")
	}

	var n big.Int
	n.Set(np)

	buf := make([]byte, 0, 8)
	var i int

	for {
		var m big.Int
		m.And(&n, big63)
		n.Rsh(&n, 6)

		buf = append(buf, CharSet64[m.Int64()])
		i++

		if n.Sign() == 0 {
			return buf, nil
		}

		n.Sub(&n, big1)
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

		n.Lsh(&n, 6)
		n.Add(&n, ind)
	}

	return n.Sub(&n, big1), nil
}

func makeBigTable(cs string) (t [256]*big.Int) {
	for i, c := range cs {
		t[c] = big.NewInt(int64(i + 1))
	}
	return t
}
