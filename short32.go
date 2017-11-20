package shortening // import "vallon.me/shortening"

const CharSet32 = `ABCDEFGHIJKLMNOPQRSTUVWXYZ234567`

var lookupTable32 = makeTable(CharSet32)

const (
	// Min values are the smallest value of 'n' that requires a string
	// of that length.
	b32min02 = 1 << 5
	b32min03 = (b32min02 + 1) << 5
	b32min04 = (b32min03 + 1) << 5
	b32min05 = (b32min04 + 1) << 5
	b32min06 = (b32min05 + 1) << 5
	b32min07 = (b32min06 + 1) << 5
	b32min08 = (b32min07 + 1) << 5
	b32min09 = (b32min08 + 1) << 5
	b32min10 = (b32min09 + 1) << 5
	b32min11 = (b32min10 + 1) << 5
	b32min12 = (b32min11 + 1) << 5
	b32min13 = (b32min12 + 1) << 5
)

var minTable32 = [...]uint64{0,
	b32min02, b32min03, b32min04, b32min05, b32min06, b32min07,
	b32min08, b32min09, b32min10, b32min11, b32min12, b32min13,
}

// Encode32 turns an uint64 into a slice of characters from 'charSet32'
func Encode32(n uint64) []byte {
	l := encLen(n, 5, minTable32[:])

	n -= b32min13

	buf := make([]byte, l)
	for i := l - 1; 0 <= i; i-- {
		buf[i], n = CharSet32[n&31], n>>5
	}

	return buf
}

// Decode32 turns a slice of characters back into the original unit64.
//
// Errors are returned for invalid characters or input that would
// cause an overflow.
func Decode32(b []byte) (n uint64, err error) {
	if 13 < len(b) || len(b) == 0 {
		return 0, InvalidDecodeLen
	}

	var invalid uint8
	for _, c := range b {
		ind := lookupTable32[c]
		invalid |= ind
		n = n<<5 | uint64(ind)
	}

	n += minTable32[len(b)-1]

	if invalid == 0xFF {
		return 0, InvalidCharacter
	}

	// b32min13 is the minimum value to have len == 13
	// any lower value is an overflow.
	if len(b) == 13 && n < b32min13 {
		return 0, Overflow
	}

	return n, nil
}
