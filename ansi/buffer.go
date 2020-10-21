package ansi

import (
	"bytes"
	"unsafe"

	"github.com/mattn/go-runewidth"
)

// Buffer is a buffer aware of ANSI escape sequences.
type Buffer struct {
	bytes.Buffer
}

// PrintableRuneWidth returns the width of all printable runes in the buffer.
func (w Buffer) PrintableRuneWidth() int {
	return PrintableRuneWidth(b2s(w.Bytes()))
}

func PrintableRuneWidth(s string) int {
	var n int
	var ansi bool

	for _, c := range s {
		if c == '\x1B' {
			// ANSI escape sequence
			ansi = true
		} else if ansi {
			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				ansi = false
			}
		} else {
			n += runewidth.RuneWidth(c)
		}
	}

	return n
}

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}
