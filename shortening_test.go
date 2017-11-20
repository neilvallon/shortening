package shortening

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

type testCoder struct {
	enc func(n uint64) []byte
	dec func(b []byte) (uint64, error)
}

func (c *testCoder) Encode(n uint64) []byte          { return c.enc(n) }
func (c *testCoder) Decode(b []byte) (uint64, error) { return c.dec(b) }

func propEncCounting(t *testing.T, charset string, enc encoder) {
	base := len(charset)
	var n uint64
	for i := -1; i < base; i++ {
		for j := 0; j < base; j++ {
			id := enc.Encode(n)
			if id[len(id)-1] != charset[j] {
				t.Fatalf("Expected rightmost char to be %q in %q", charset[j], id)
			}

			if i != -1 && id[0] != charset[i] {
				t.Fatalf("Expected leftmost char to be %q in %q", charset[i], id)
			}

			n++
		}
	}
}

func propRollover(t *testing.T, charset string, table []uint64, enc encoder) {
	firstChar := charset[0]
	lastChar := charset[len(charset)-1]

	for i := 0; i < len(table); i++ {
		a := enc.Encode(table[i] - 1)
		b := enc.Encode(table[i])

		if len(a)+1 != len(b) {
			t.Fatal("error")
		}

		for _, c := range a {
			if c != lastChar {
				t.Fatalf("expected all chars to be %q in %q", lastChar, a)
			}
		}

		for _, c := range b {
			if c != firstChar {
				t.Fatalf("expected all chars to be %q in %q", firstChar, b)
			}
		}
	}
}

func propDecCounting(t *testing.T, charset string, dec decoder) {
	base := len(charset)
	var n uint64
	for i := -1; i < base; i++ {
		for j := 0; j < base; j++ {
			id := []byte{charset[j]}
			if i != -1 {
				id = []byte{charset[i], charset[j]}
			}

			nn, err := dec.Decode(id)
			if err != nil {
				t.Logf("%q", id)
				t.Fatal(err)
			}

			if nn != n {
				t.Errorf("got: %d - expected: %d", nn, n)
			}

			n++
		}
	}
}

func propParity(t *testing.T, c coder) {
	// test first 10k
	for i := uint64(0); i < 100000; i++ {
		n, err := c.Decode(c.Encode(i))
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
		n, err := c.Decode(c.Encode(i))
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

		n, err := c.Decode(c.Encode(v))
		if err != nil {
			t.Logf("test ID: %d", v)
			t.Fatal(err)
		}
		if n != v {
			t.Errorf("got: %d - expected: %d", n, v)
		}
	}
}

const benchSet = 1 << 13 // 8192

func benchEncode(b *testing.B, enc encoder) {
	ids := make([]uint64, benchSet)
	for i := range ids {
		ids[i] = rand.Uint64()
	}

	b.Run("encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = enc.Encode(ids[i&(benchSet-1)])
		}
	})
}

func benchDecode(b *testing.B, c coder) {
	urls := make([][]byte, benchSet)
	for i := range urls {
		urls[i] = c.Encode(rand.Uint64())
	}

	b.Run("decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = c.Decode(urls[i&(benchSet-1)])
		}
	})
}
