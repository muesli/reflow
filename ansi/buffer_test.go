package ansi

import (
	"bytes"
	"testing"
)

func TestBuffer_PrintableRuneWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    string
		len  int
	}{
		{"CSI", "\x1B[38;2;249;38;114mfoo", 3},
		// sixel example from wikipedia
		{"DCS", "\x1BPq#0;2;0;0;0#1;2;100;100;0#2;2;0;100;0#1~~@@vv@@~~@@~~$#2??}}GG}}??}}??-#1!14@\x1B\\foo", 3},
		{"wide_chars", "Hello, 世界", 11},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			var bb bytes.Buffer
			bb.WriteString(tst.s)
			b := Buffer{bb}

			if n := b.PrintableRuneWidth(); n != tst.len {
				t.Fatalf("width should be %d, got %d", tst.len, n)
			}
		})
	}
}

// go test -bench=Benchmark_PrintableRuneWidth -benchmem -count=4
func Benchmark_PrintableRuneWidth(b *testing.B) {
	s := "\x1B[38;2;249;38;114mfoo"

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			PrintableRuneWidth(s)
		}
	})
}
