package ansi

import (
	"bytes"

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

// Truncate truncates a string at the given printable cell width, leaving any
// ansi sequences intact.
func Truncate(s string, w int) string {
	return TruncateWithTail(s, w, "")
}

// TruncateWithTail truncates a string at the given printable cell width,
// leaving any ansi sequences intact. A tail is then added to the end of the
// string.
func TruncateWithTail(s string, w int, tail string) string {
	if PrintableRuneWidth(s) <= w {
		return s
	}

	const ansiReset = "\x1B[0m"

	if tail != "" {
		tail += ansiReset
	}

	tw := PrintableRuneWidth(tail)
	w -= tw
	if w < 0 {
		return tail
	}

	r := []rune(s)
	ansi := false
	n := 0
	i := 0

	for ; i < len(r); i++ {
		accPrintableRuneWidth(r[i], &n, &ansi)
		if n > w {
			break
		}
	}

	return string(r[0:i]) + ansiReset + tail
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
