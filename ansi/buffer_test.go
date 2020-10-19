package ansi

import (
	"bytes"
	"testing"
)

func TestBuffer_PrintableRuneWidth(t *testing.T) {
	var bb bytes.Buffer
	bb.WriteString("\x1B[38;2;249;38;114mfoo")
	b := Buffer{bb}

	if n := b.PrintableRuneWidth(); n != 3 {
		t.Fatalf("width should be 3, got %d", n)
	}
}

func BenchmarkPrintableRuneWidth(b *testing.B) {
	s := "\x1B[38;2;249;38;114mfoo"
	var n int

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n = PrintableRuneWidth(s)
	}

	if n != 3 {
		b.Fatalf("width should be 3, got %d", n)
	}
}
