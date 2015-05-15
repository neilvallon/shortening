package shortening

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestEncodeShareInt(t *testing.T) {
	tests := []struct {
		ID  uint64
		URL string
	}{
		{0, "A"},
		{1, "B"},
		{2, "C"},
		{63, "_"},
		{64, "AA"},
		{65, "BA"},
		{66, "CA"},
		{128, "AB"},
		{128, "AB"},
		{4159, "__"},
		{4160, "AAA"},
		{1090785344, "AAAAAA"},

		{1171221845949812800 - 1, "__________"},
		{1171221845949812800, "AAAAAAAAAAA"},
		{1171221845949812800 + 1, "BAAAAAAAAAA"},

		{math.MaxUint64 - 1, "----------O"},
		{math.MaxUint64, "_---------O"},
	}

	for _, test := range tests {
		if e := string(Encode(test.ID)); e != test.URL {
			t.Errorf("got: %q - expected: %q", e, test.URL)
		}
	}
}

func TestDecodeShareInt(t *testing.T) {
	tests := []struct {
		ID  uint64
		URL string
	}{
		{0, "A"},
		{1, "B"},
		{2, "C"},
		{63, "_"},
		{64, "AA"},
		{65, "BA"},
		{66, "CA"},
		{128, "AB"},
		{128, "AB"},
		{4159, "__"},
		{4160, "AAA"},
		{1090785344, "AAAAAA"},

		{1171221845949812800 - 1, "__________"},
		{1171221845949812800, "AAAAAAAAAAA"},
		{1171221845949812800 + 1, "BAAAAAAAAAA"},

		{math.MaxUint64 - 1, "----------O"},
		{math.MaxUint64, "_---------O"},
	}

	for _, test := range tests {
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
		v := uint64(rnd.Uint32())

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
		ids[i] = uint64(rand.Int63())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode(ids[i%benchSet])
	}
}

func BenchmarkDecode(b *testing.B) {
	urls := make([][]byte, benchSet)
	for i := range urls {
		urls[i] = Encode(uint64(rand.Int63()))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(urls[i%benchSet])
	}
}
