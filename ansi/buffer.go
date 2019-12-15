package ansi

import (
	"bytes"

	"github.com/mattn/go-runewidth"
)

// Buffer is a buffer aware of ANSI escape sequences.
type Buffer struct {
	bytes.Buffer
}

// PrintableRuneCount returns the amount of printable runes in the buffer.
func (w Buffer) PrintableRuneCount() int {
	var n int
	var ansi bool
	for _, c := range w.String() {
		if c == '\x1B' {
			// ANSI escape sequence
			ansi = true
		} else if ansi {
			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				ansi = false
			}
		} else {
			n += runewidth.StringWidth(string(c))
		}
	}

	return n
}
