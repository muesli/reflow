package ansi

import (
	"bytes"

	"github.com/rivo/uniseg"
)

// Buffer is a buffer aware of ANSI escape sequences.
type Buffer struct {
	bytes.Buffer
}

// PrintableRuneWidth returns the cell width of all printable runes in the
// buffer.
func (w Buffer) PrintableRuneWidth() int {
	return PrintableRuneWidth(w.String())
}

// PrintableRuneWidth returns the cell width of the given string.
func PrintableRuneWidth(s string) int {
	n := make([]rune, 0, len(s))
	var ansi bool

	for _, c := range s {
		switch {
		case c == Marker:
			// ANSI escape sequence
			ansi = true
		case ansi && IsTerminator(c):
			// ANSI sequence terminated
			ansi = false
		case ansi:
		default:
			n = append(n, c)
		}
	}

	return uniseg.StringWidth(string(n))
}
