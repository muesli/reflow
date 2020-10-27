package ansi

import (
	"bytes"
	"fmt"
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

	tests := []struct {
		in            string
		out           string
		width         int
		expectedWidth int
	}{
		{
			"\x1B[38;2;249;38;114m你\x1B[7m好\x1B[0m",
			"\x1B[38;2;249;38;114m你\x1B[7m\x1B[0m",
			2,
			2,
		},
		{
			"\x1B[38;2;249;38;114m你\x1B[7m好\x1B[0m",
			"\x1B[38;2;249;38;114m\x1B[0m",
			1,
			0,
		},
		{
			"It’s me!",
			"It’s me!",
			10,
			8,
		},
		{
			"It’s \x1B[7mme!",
			"It’s \x1B[7m\x1B[0m",
			5,
			5,
		},
	}

	i := 0
	for _, tt := range tests {
		t.Run(fmt.Sprintf("truncate-%d", i), func(t *testing.T) {
			t.Parallel()
			res := Truncate(tt.in, tt.width)
			if n := PrintableRuneWidth(res); n != tt.expectedWidth {
				t.Fatalf("width should be %d, got %d", tt.expectedWidth, n)
			}
			if res != tt.out {
				t.Fatalf("expected '%s' got '%s'\x1B[0m", tt.out, res)
			}
		})
		i++
	}
}

// go test -bench=Benchmark_Truncate -benchmem -count=4
func Benchmark_Truncate(b *testing.B) {
	s := "\x1B[38;2;249;38;114m你\x1B[7m好\x1B[0m"

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			Truncate(s, 2)
		}
	})
}
