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

func BenchmarkEncode(b *testing.B) {
	b.StopTimer()
	ids := make([]uint64, b.N)
	for i := range ids {
		ids[i] = uint64(rand.Int63())
	}

	b.StartTimer()
	for _, id := range ids {
		_ = Encode(id)
	}
}

func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	urls := make([][]byte, b.N)
	for i := range urls {
		urls[i] = Encode(uint64(rand.Int63()))
	}

	b.StartTimer()
	for _, b := range urls {
		_, _ = Decode(b)
	}
}