package shortening

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var tests32 = []struct {
	ID  uint64
	URL string
}{
	{0, "A"},
	{1, "B"},
	{2, "C"},
	{b32min02 - 1, "7"},
	{b32min02, "AA"},
	{b32min02 + 1, "BA"},
	{b32min02 + 2, "CA"},
	{b32min02 * 2, "AB"},
	{b32min03 - 1, "77"},
	{b32min03, "AAA"},
	{b32min06, "AAAAAA"},

	{b32min13 - 1, "777777777777"},
	{b32min13, "AAAAAAAAAAAAA"},
	{b32min13 + 1, "BAAAAAAAAAAAA"},

	{math.MaxUint64 - 1, "666666666666O"},
	{math.MaxUint64, "766666666666O"},
}

func TestEncode32ShareInt(t *testing.T) {
	for _, test := range tests32 {
		if e := string(Encode32(test.ID)); e != test.URL {
			t.Errorf("got: %q - expected: %q", e, test.URL)
		}
	}
}

func TestDecode32ShareInt(t *testing.T) {
	for _, test := range tests32 {
		n, err := Decode32([]byte(test.URL))
		if err != nil {
			t.Logf("%q", test.URL)
			t.Fatal(err)
		}
		if n != test.ID {
			t.Errorf("got: %d - expected: %d", n, test.ID)
		}
	}
}

func TestDecode32Errors(t *testing.T) {
	tests := [][]byte{
		nil,
		[]byte(""),
		[]byte("A76666666666O"),  // overflow
		[]byte("AAAAAAAAAAAAAA"), // 14+ bytes
		[]byte("*"),              // invalid character
		[]byte("\xFF"),           // invalid character
	}

	for _, test := range tests {
		if _, err := Decode32(test); err == nil {
			t.Logf("%q", test)
			t.Error("error expected. got nil")
		}
	}
}

func TestEncDec32Parity(t *testing.T) {
	// test first 10k
	for i := uint64(0); i < 100000; i++ {
		n, err := Decode32(Encode32(i))
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
		n, err := Decode32(Encode32(i))
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
		v := uint64(rnd.Uint64())

		n, err := Decode32(Encode32(v))
		if err != nil {
			t.Logf("test ID: %d", v)
			t.Fatal(err)
		}
		if n != v {
			t.Errorf("got: %d - expected: %d", n, v)
		}
	}
}

func BenchmarkEncode32(b *testing.B) {
	ids := make([]uint64, benchSet)
	for i := range ids {
		ids[i] = rand.Uint64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode32(ids[i%benchSet])
	}
}

func BenchmarkDecode32(b *testing.B) {
	urls := make([][]byte, benchSet)
	for i := range urls {
		urls[i] = Encode32(rand.Uint64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode32(urls[i%benchSet])
	}
}
