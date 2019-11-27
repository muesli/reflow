package reflow

import (
	"bytes"
	"unicode"
)

var (
	defaultBreakpoints = []rune{'-'}
	defaultNewline     = []rune{'\n'}
)

type Reflow struct {
	Limit       int
	Breakpoints []rune
	Newline     []rune

	buf   bytes.Buffer
	word  bytes.Buffer
	space bytes.Buffer

	lineLen int
	ansi    bool
}

func NewReflow(limit int) *Reflow {
	return &Reflow{
		Limit:       limit,
		Breakpoints: defaultBreakpoints,
		Newline:     defaultNewline,
	}
}

func ReflowBytes(b []byte, limit int) []byte {
	f := NewReflow(limit)
	f.Write(b)
	f.Close()

	return f.Bytes()
}

func ReflowString(s string, limit int) string {
	return string(ReflowBytes([]byte(s), limit))
}

func (w *Reflow) addSpace() {
	w.lineLen += w.space.Len()
	w.buf.Write(w.space.Bytes())
	w.space.Reset()
}

func (w *Reflow) addWord() {
	if w.word.Len() > 0 {
		w.addSpace()
		w.lineLen += w.word.Len()
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

func (w *Reflow) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == '\x1B' {
			// ANSI escape sequence
			w.addSpace()
			w.addWord()

			w.buf.WriteRune(c)
			w.ansi = true
		} else if w.ansi {
			w.buf.WriteRune(c)
			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				w.ansi = false
			}
		} else if inGroup(w.Newline, c) {
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
			w.addWord()
			w.buf.WriteRune(c)
		} else {
			// any other character
			w.word.WriteRune(c)

			// add a line break if the current word would exceed the line's
			// character limit
			if w.lineLen+w.space.Len()+w.word.Len() > w.Limit && w.word.Len() < w.Limit {
				w.addNewLine()
			}
		}
	}

	return len(b), nil
}

func (w *Reflow) Close() error {
	w.addWord()
	return nil
}

func (w *Reflow) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Reflow) String() string {
	return w.buf.String()
}
