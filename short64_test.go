package shortening

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var b64tests = []struct {
	ID  uint64
	URL string
}{
	{0, "A"},
	{1, "B"},
	{2, "C"},
	{min02 - 1, "_"},
	{min02, "AA"},
	{min02 + 1, "BA"},
	{min02 + 2, "CA"},
	{min02 * 2, "AB"},
	{min03 - 1, "__"},
	{min03, "AAA"},
	{min06, "AAAAAA"},

	{min11 - 1, "__________"},
	{min11, "AAAAAAAAAAA"},
	{min11 + 1, "BAAAAAAAAAA"},

	{math.MaxUint64 - 1, "----------O"},
	{math.MaxUint64, "_---------O"},
}

func TestEncodeShareInt(t *testing.T) {
	for _, test := range b64tests {
		if e := string(Encode(test.ID)); e != test.URL {
			t.Errorf("got: %q - expected: %q", e, test.URL)
		}
	}
}

func TestDecodeShareInt(t *testing.T) {
	for _, test := range b64tests {
		n, err := Decode([]byte(test.URL))
		if err != nil {
			t.Logf("%q", test.URL)
			t.Fatal(err)
		}
		if n != test.ID {
			t.Errorf("got: %d - expected: %d", n, test.ID)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	tests := [][]byte{
		nil,
		[]byte(""),
		[]byte("123456789_-"),  // overflow
		[]byte("AAAAAAAAAAAA"), // 12+ bytes
		[]byte("*"),            // invalid character
		[]byte("\xFF"),         // invalid character
	}

	for _, test := range tests {
		if _, err := Decode(test); err == nil {
			t.Logf("%q", test)
			t.Error("error expected. got nil")
		}
	}
}

func TestEncDecParity(t *testing.T) {
	// test first 10k
	for i := uint64(0); i < 100000; i++ {
		n, err := Decode(Encode(i))
		if err != nil {
			t.Logf("test ID: %d", i)
			t.Fatal(err)
		}
		if n != i {
			t.Errorf("got: %d - expected: %d", n, i)
		}
	}

	// test last 10k
	for i := uint64(math.MaxUint64 - 100000); i < math.MaxUint64; i++ {
		n, err := Decode(Encode(i))
		if err != nil {
			t.Logf("test ID: %d", i)
			t.Fatal(err)
		}
		if n != i {
			t.Errorf("got: %d - expected: %d", n, i)
		}
	}

	// test random 10k
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100000; i++ {
		v := rnd.Uint64()

		n, err := Decode(Encode(v))
		if err != nil {
			t.Logf("test ID: %d", v)
			t.Fatal(err)
		}
		if n != v {
			t.Errorf("got: %d - expected: %d", n, v)
		}
	}
}

const benchSet = 10000

func BenchmarkEncode(b *testing.B) {
	ids := make([]uint64, benchSet)
	for i := range ids {
		ids[i] = rand.Uint64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode(ids[i%benchSet])
	}
}

func BenchmarkDecode(b *testing.B) {
	urls := make([][]byte, benchSet)
	for i := range urls {
		urls[i] = Encode(rand.Uint64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(urls[i%benchSet])
	}
}
