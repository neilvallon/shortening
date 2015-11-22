package shortening

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestEncBigInt(t *testing.T) {
	tests := []struct {
		ID  int64
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
	}

	for _, test := range tests {
		s, err := EncodeBig(big.NewInt(test.ID))
		if err != nil {
			t.Error(err)
		}

		if e := string(s); e != test.URL {
			t.Errorf("got: %q - expected: %q", e, test.URL)
		}
	}
}

func TestDecBigInt(t *testing.T) {
	tests := []struct {
		ID  int64
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
	}

	for _, test := range tests {
		n, err := DecodeBig([]byte(test.URL))
		if err != nil {
			t.Logf("%q", test.URL)
			t.Fatal(err)
		}
		if n.Int64() != test.ID {
			t.Errorf("got: %d - expected: %d", n, test.ID)
		}
	}
}

func TestBigDecodeErrors(t *testing.T) {
	tests := [][]byte{
		nil,
		[]byte(""),
		[]byte("*"),    // invalid character
		[]byte("\xFF"), // invalid character
	}

	for _, test := range tests {
		if _, err := DecodeBig(test); err == nil {
			t.Logf("%q", test)
			t.Error("error expected. got nil")
		}
	}
}

func TestBigEncDecParity(t *testing.T) {
	// test first 10k
	for i := int64(0); i < 100000; i++ {
		s, _ := EncodeBig(big.NewInt(i))
		n, err := DecodeBig(s)

		if err != nil {
			t.Logf("test ID: %d", i)
			t.Fatal(err)
		}
		if n.Int64() != i {
			t.Errorf("got: %d - expected: %d", n, i)
		}
	}

	// test random 10k
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100000; i++ {
		v := int64(rnd.Int())
		s, _ := EncodeBig(big.NewInt(v))

		n, err := DecodeBig(s)
		if err != nil {
			t.Logf("test ID: %d", v)
			t.Fatal(err)
		}
		if n.Int64() != v {
			t.Errorf("got: %d - expected: %d", n, v)
		}
	}
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func BenchmarkBigEncode(b *testing.B) {
	ids := make([]big.Int, benchSet)
	for i := range ids {
		var bi big.Int
		bi.Rand(rng, big.NewInt(math.MaxInt64))
		ids[i] = bi
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EncodeBig(&ids[i%benchSet])
	}
}

func BenchmarkBigDecode(b *testing.B) {
	urls := make([][]byte, benchSet)
	for i := range urls {
		urls[i] = Encode(uint64(rand.Int63()))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeBig(urls[i%benchSet])
	}
}
