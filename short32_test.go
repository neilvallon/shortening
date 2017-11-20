package shortening

import (
	"math"
	"testing"
)

var testCoder32 = &testCoder{Encode32, Decode32}

func TestEncode32(t *testing.T) {
	// verify counting for two place values
	propEncCounting(t, CharSet32, testCoder32)

	// checks expected values before and after a place value change.
	propRollover(t, CharSet32, minTable32[1:], testCoder32)

	// manualy check max uint64 encode
	maxEncoded := "O666666666667"
	if e := string(testCoder32.Encode(math.MaxUint64)); e != maxEncoded {
		t.Errorf("got: %q - expected: %q", e, maxEncoded)
	}
}

func TestDecode32(t *testing.T) {
	propDecCounting(t, CharSet32, testCoder32)

	// manualy check max uint64 encode
	n, err := testCoder32.Decode([]byte("O666666666667"))
	if err != nil {
		t.Error(err)
	}

	if max := uint64(math.MaxUint64); n != max {
		t.Errorf("got: %q - expected: %q", n, max)
	}
}

func TestDecode32Errors(t *testing.T) {
	tests := []struct {
		ID  []byte
		Err error
	}{
		{nil, InvalidDecodeLen},
		{[]byte(""), InvalidDecodeLen},
		{[]byte("O66666666667A"), Overflow},
		{[]byte("AAAAAAAAAAAAAA"), InvalidDecodeLen},
		{[]byte("*"), InvalidCharacter},
		{[]byte("\xFF"), InvalidCharacter},
	}

	for _, test := range tests {
		_, err := testCoder32.Decode(test.ID)
		if err == nil {
			t.Logf("%q", test.ID)
			t.Error("error expected. got nil")
		} else if err != test.Err {
			t.Logf("%q", test)
			t.Fatalf("got: %q - expected: %q", err, test.Err)
		}
	}
}

func TestEncDec32Parity(t *testing.T) { propParity(t, testCoder32) }

func BenchmarkEncode32(b *testing.B) { benchEncode(b, testCoder32) }
func BenchmarkDecode32(b *testing.B) { benchDecode(b, testCoder32) }
