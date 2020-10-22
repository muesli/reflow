package ansi

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
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
	var n int
	var ansi bool

	for _, c := range s {
		accPrintableRuneWidth(c, &n, &ansi)
	}

	return n
}

// Truncate truncates a given string at the given printable cell width, leaving
// any ansi sequences intact.
func Truncate(s string, w int, tail string) string {
	if PrintableRuneWidth(s) <= w {
		return s
	}

	var n int
	var ansi bool
	var acc strings.Builder

	for _, c := range s {
		accPrintableRuneWidth(c, &n, &ansi)
		if n > w {
			break
		}

		_, _ = acc.WriteRune(c)
	}

	_, _ = acc.WriteString("\x1B[0m") // terminate any open ANSI sequences
	_, _ = acc.WriteString(tail)
	return acc.String()
}

// Used to accumulate the printable rune width while tracking whether we're in
// an ansi sequence.
func accPrintableRuneWidth(c rune, n *int, ansi *bool) {
	if c == '\x1B' {
		// ANSI escape sequence
		*ansi = true
	} else if *ansi {
		if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
			// ANSI sequence terminated
			*ansi = false
		}
	} else {
		*n += runewidth.RuneWidth(c)
	}
}
