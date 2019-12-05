package reflow

import (
	"bytes"
	"unicode"
)

var (
	defaultBreakpoints = []rune{'-'}
	defaultNewline     = []rune{'\n'}
)

// Reflow contains settings and state for customisable text reflowing with
// support for ANSI escape sequences. This means you can style your terminal
// output without affecting the word wrapping algorithm.
type Reflow struct {
	Limit        int
	Breakpoints  []rune
	Newline      []rune
	KeepNewlines bool

	buf   bytes.Buffer
	space bytes.Buffer
	word  ANSIBuffer

	lineLen int
	ansi    bool
}

// NewReflow returns a new instance of Reflow, initialized with defaults.
func NewReflow(limit int) *Reflow {
	return &Reflow{
		Limit:        limit,
		Breakpoints:  defaultBreakpoints,
		Newline:      defaultNewline,
		KeepNewlines: true,
	}
}

// Bytes is shorthand for declaring a new default Reflow instance,
// used to immediately reflow a byte slice.
func Bytes(b []byte, limit int) []byte {
	f := NewReflow(limit)
	_, _ = f.Write(b)
	f.Close()

	return f.Bytes()
}

// String is shorthand for declaring a new default Reflow instance,
// used to immediately reflow a string.
func String(s string, limit int) string {
	return string(Bytes([]byte(s), limit))
}

func (w *Reflow) addSpace() {
	w.lineLen += w.space.Len()
	w.buf.Write(w.space.Bytes())
	w.space.Reset()
}

func (w *Reflow) addWord() {
	if w.word.Len() > 0 {
		w.addSpace()
		w.lineLen += w.word.PrintableRuneCount()
		w.buf.Write(w.word.Bytes())
		w.word.Reset()
	}
}

func (w *Reflow) addNewLine() {
	w.buf.WriteRune('\n')
	w.lineLen = 0
	w.space.Reset()
}

func inGroup(a []rune, c rune) bool {
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}

// Write is used to write more content to the reflow buffer.
func (w *Reflow) Write(b []byte) (int, error) {
	if w.Limit == 0 {
		return w.buf.Write(b)
	}

	for _, c := range string(b) {
		if c == '\x1B' {
			// ANSI escape sequence
			w.word.WriteRune(c)
			w.ansi = true
		} else if w.ansi {
			w.word.WriteRune(c)
			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				w.ansi = false
			}
		} else if w.KeepNewlines && inGroup(w.Newline, c) {
			// end of current line
			// see if we can add the content of the space buffer to the current line
			if w.word.Len() == 0 {
				if w.lineLen+w.space.Len() > w.Limit {
					w.lineLen = 0
				} else {
					// preserve whitespace
					w.buf.Write(w.space.Bytes())
				}
				w.space.Reset()
			}

			w.addWord()
			w.addNewLine()
		} else if unicode.IsSpace(c) {
			// end of current word
			w.addWord()
			w.space.WriteRune(c)
		} else if inGroup(w.Breakpoints, c) {
			// valid breakpoint
			w.addSpace()
			w.addWord()
			w.buf.WriteRune(c)
		} else {
			// any other character
			w.word.WriteRune(c)

			// add a line break if the current word would exceed the line's
			// character limit
			if w.lineLen+w.space.Len()+w.word.PrintableRuneCount() > w.Limit &&
				w.word.PrintableRuneCount() < w.Limit {
				w.addNewLine()
			}
		}
	}

	return len(b), nil
}

// Close will finish the reflow operation. Always call it before trying to
// retrieve the final result.
func (w *Reflow) Close() error {
	w.addWord()
	return nil
}

// Bytes returns the reflow result as a byte slice.
func (w *Reflow) Bytes() []byte {
	return w.buf.Bytes()
}

// String returns the reflow result as a string.
func (w *Reflow) String() string {
	return w.buf.String()
}
