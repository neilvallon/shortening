package shortening // import "vallon.me/shortening"

import "errors"

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

// Encode turns an uint64 into a slice of characters from 'charSet'
func Encode(n uint64) []byte {
	var buf [11]byte
	start := 10

	switch {
	case min11 <= n:
		buf[0], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min10 <= n:
		buf[1], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min09 <= n:
		buf[2], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min08 <= n:
		buf[3], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min07 <= n:
		buf[4], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min06 <= n:
		buf[5], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min05 <= n:
		buf[6], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min04 <= n:
		buf[7], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min03 <= n:
		buf[8], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	case min02 <= n:
		buf[9], n = charSet[n&63], (n>>6)-1
		start--
		fallthrough
	default:
		buf[10] = charSet[n&63]
	}

	return buf[start:]
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
		n = n<<6 + uint64(ind)
		fallthrough
	case 10:
		ind = lookupTable[b[9]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 9:
		ind = lookupTable[b[8]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 8:
		ind = lookupTable[b[7]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 7:
		ind = lookupTable[b[6]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 6:
		ind = lookupTable[b[5]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 5:
		ind = lookupTable[b[4]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 4:
		ind = lookupTable[b[3]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 3:
		ind = lookupTable[b[2]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 2:
		ind = lookupTable[b[1]]
		invalid |= ind
		n = n<<6 + uint64(ind)
		fallthrough
	case 1:
		ind = lookupTable[b[0]]
		invalid |= ind
		n = n<<6 + uint64(ind)
	}

	if invalid == 0xFF {
		return 0, errors.New("shortening: invalid decode character")
	}

	// 1171221845949812800 is the minimum value to have len == 11
	// any lower value is an overflow.
	if len(b) == 11 && n-1 < min11 {
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
