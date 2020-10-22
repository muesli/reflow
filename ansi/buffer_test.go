package ansi

import (
	"bytes"
	"testing"
)

func TestBuffer_PrintableRuneWidth(t *testing.T) {
	t.Parallel()

	var bb bytes.Buffer
	bb.WriteString("\x1B[38;2;249;38;114mfoo")
	b := Buffer{bb}

	if n := b.PrintableRuneWidth(); n != 3 {
		t.Fatalf("width should be 3, got %d", n)
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

func Test_Truncate(t *testing.T) {
	t.Parallel()

	s := "\x1B[38;2;249;38;114m你\x1B[7m好\x1B[0m"

	if n := PrintableRuneWidth(Truncate(s, 2, "")); n != 2 {
		t.Fatalf("width should be 2, got %d", n)
	}
}

// go test -bench=Benchmark_Truncate -benchmem -count=4
func Benchmark_Truncate(b *testing.B) {
	s := "\x1B[38;2;249;38;114m你\x1B[7m好\x1B[0m"

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			Truncate(s, 2, "")
		}
	})
}
