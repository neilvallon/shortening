package shortening

import (
	"math"
	"testing"
)

var testCoder64 = &testCoder{Encode, Decode}

func TestEncode64(t *testing.T) {
	// verify counting for two place values
	propEncCounting(t, CharSet64, testCoder64)

	// checks expected values before and after a place value change.
	propRollover(t, CharSet64, minTable[1:], testCoder64)

	// manualy check max uint64 encode
	maxEncoded := "O---------_"
	if e := string(testCoder64.Encode(math.MaxUint64)); e != maxEncoded {
		t.Errorf("got: %q - expected: %q", e, maxEncoded)
	}
}

func TestDecode64(t *testing.T) {
	propDecCounting(t, CharSet64, testCoder64)

	// manualy check max uint64 encode
	n, err := testCoder64.Decode([]byte("O---------_"))
	if err != nil {
		t.Error(err)
	}

	if max := uint64(math.MaxUint64); n != max {
		t.Errorf("got: %q - expected: %q", n, max)
	}
}

func TestDecode64Errors(t *testing.T) {
	tests := []struct {
		ID  []byte
		Err error
	}{
		{nil, InvalidDecodeLen},
		{[]byte(""), InvalidDecodeLen},
		{[]byte("-_987654321"), Overflow},
		{[]byte("AAAAAAAAAAAA"), InvalidDecodeLen},
		{[]byte("*"), InvalidCharacter},
		{[]byte("\xFF"), InvalidCharacter},
	}

	for _, test := range tests {
		_, err := testCoder64.Decode(test.ID)
		if err == nil {
			t.Logf("%q", test.ID)
			t.Error("error expected. got nil")
		} else if err != test.Err {
			t.Logf("%q", test)
			t.Fatalf("got: %q - expected: %q", err, test.Err)
		}
	}
}

func TestEncDec64Parity(t *testing.T) { propParity(t, testCoder64) }

func BenchmarkEncode64(b *testing.B) { benchEncode(b, testCoder64) }
func BenchmarkDecode64(b *testing.B) { benchDecode(b, testCoder64) }
