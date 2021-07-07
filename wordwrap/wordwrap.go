package wordwrap

import (
	"bytes"
	"strings"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

var (
	defaultBreakpoints = []rune{'-'}
	defaultNewline     = []rune{'\n'}
)

// WordWrap contains settings and state for customisable text reflowing with
// support for ANSI escape sequences. This means you can style your terminal
// output without affecting the word wrapping algorithm.
type WordWrap struct {
	Limit          int
	Breakpoints    []rune
	Newline        []rune
	KeepNewlines   bool
	HardWrap       bool
	TabReplace     string // since tabs can have different lengths, replace them with this when hardwrap is enabled
	PreserveSpaces bool

	buf   bytes.Buffer // processed and, in line, accepted bytes
	space bytes.Buffer // pending continues spaces bytes
	word  ansi.Buffer  // pending continues word bytes

	lineLen int // the visible length of the line not accurate for tabs
	ansi    bool

	wroteBegin bool         // mark is since the last newline something has written to the buffer (for ansi restart)
	lastAnsi   bytes.Buffer // hold last active ansi sequence

	// the following are used to remove leading zeros from the single arguments of the ansi-sequence, but still detect single zeros:
	// \x1B[0031;0000m => \x1B[31;0m
	newArgument bool
	leadingZero bool
}

// NewWriter returns a new instance of a word-wrapping writer, initialized with
// default settings.
func NewWriter(limit int) *WordWrap {
	return &WordWrap{
		Limit:        limit,
		Breakpoints:  defaultBreakpoints,
		Newline:      defaultNewline,
		KeepNewlines: true,
	}
}

// Bytes is shorthand for declaring a new default WordWrap instance,
// used to immediately word-wrap a byte slice.
func Bytes(b []byte, limit int) []byte {
	f := NewWriter(limit)
	_, _ = f.Write(b)
	_ = f.Close()

	return f.Bytes()
}

// String is shorthand for declaring a new default WordWrap instance,
// used to immediately word-wrap a string.
func String(s string, limit int) string {
	return string(Bytes([]byte(s), limit))
}

// HardWrap is a shorthand for declaring a new hardwrapping WordWrap instance,
// since variable length characters can not be hard wrapped to a fixed length,
// tabs will be replaced by TabReplace, use according amount of spaces.
func HardWrap(s string, limit int, tabReplace string) string {
	f := NewWriter(limit)
	f.HardWrap = true
	f.TabReplace = tabReplace
	_, _ = f.Write([]byte(s))
	f.Close()

	return f.String()
}

// adds pending spaces to the buf(fer) and then resets the space buffer.
func (w *WordWrap) addSpace() {
	if w.space.Len() <= w.Limit-w.lineLen {
		w.lineLen += w.space.Len()
		_, _ = w.buf.Write(w.space.Bytes())
	} else {
		length := w.space.Len()
		first := w.Limit - w.lineLen
		_, _ = w.buf.WriteString(strings.Repeat(" ", first))
		length -= first
		for length >= w.Limit {
			_, _ = w.buf.WriteString("\n" + strings.Repeat(" ", w.Limit))
			length -= w.Limit
		}
		if length > 0 {
			_, _ = w.buf.WriteString("\n" + strings.Repeat(" ", length))
		}
		w.lineLen = length
	}
	w.space.Reset()
}

func (w *WordWrap) addWord() {
	if w.word.Len() > 0 {
		w.addSpace()
		w.lineLen += w.word.PrintableRuneWidth()
		_, _ = w.buf.Write(w.word.Bytes())
		w.word.Reset()
	}
}

func (w *WordWrap) addNewLine() {
	if w.PreserveSpaces {
		w.addSpace()
	}
	if w.lastAnsi.Len() != 0 {
		// end ansi before linebreak
		_, _ = w.buf.WriteString("\x1B[0m")
	}
	_, _ = w.buf.WriteRune('\n')
	w.lineLen = 0
	w.space.Reset()
	w.wroteBegin = false
}

func inGroup(a []rune, c rune) bool {
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}

// Write is used to write more content to the word-wrap buffer.
func (w *WordWrap) Write(b []byte) (int, error) {
	if w.Limit == 0 {
		return w.buf.Write(b)
	}

	s := string(b)
	if !w.KeepNewlines {
		s = strings.Replace(strings.TrimSpace(s), "\n", " ", -1)
	}

	if w.HardWrap {
		s = strings.Replace(s, "\t", w.TabReplace, -1)
	}

	for _, c := range s {
		// Restart Ansi after line break if there is more text
		if !w.wroteBegin && !w.ansi && w.lastAnsi.Len() != 0 {
			_, _ = w.buf.Write(w.lastAnsi.Bytes())
			w.addWord()
		}
		w.wroteBegin = true
		if c == '\x1B' {
			// ANSI escape sequence
			_, _ = w.word.WriteRune(c)
			_, _ = w.lastAnsi.WriteRune(c)
			w.ansi = true
			w.newArgument = true
		} else if w.ansi {

			// ignore leading zeros but remember single ones.
			if c == '0' && w.newArgument {
				w.leadingZero = true
				continue
			}
			w.newArgument = false
			// if a digit other then zero is encountered reset leading zero since we can ignore the leading zeroes if there where any.
			if inGroup([]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}, c) {
				w.leadingZero = false
			}

			// check if new ANSI-argument starts
			if inGroup([]rune{'[', ';'}, c) {
				w.newArgument = true
				// if w.leadingZero is here, we know that its a valid zero => reset and restart sequence.
				if w.leadingZero {
					// since we are still in the middle of the sequence and have reset the last ansi, we have to restart a new sequence:
					w.lastAnsi.Reset()
					_, _ = w.lastAnsi.WriteString("\x1B[")
					w.leadingZero = false
					_, _ = w.word.WriteString("0m\x1B[")
					// "\x1B[31;0;32m" => "\x1B[31;0m\x1B[32m"
					continue // dont write "replace" semicolon
				}
			}

			_, _ = w.lastAnsi.WriteRune(c)

			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// dont restart lastAnsi since its a end of a sequence. (not in the middle of one)
				if w.leadingZero {
					_, _ = w.word.WriteRune('0')

					w.lastAnsi.Reset()
					w.leadingZero = false
				}
				// ANSI sequence terminated
				w.ansi = false
			}

			_, _ = w.word.WriteRune(c)

		} else if inGroup(w.Newline, c) {
			// end of current line
			// see if we can add the content of the space buffer to the current line
			if w.word.Len() == 0 {
				if w.lineLen+w.space.Len() > w.Limit {
					w.lineLen = 0
				} else {
					// preserve whitespace
					_, _ = w.buf.Write(w.space.Bytes())
				}
				w.space.Reset()
			}

			w.addWord()
			w.addNewLine()
		} else if unicode.IsSpace(c) {
			// end of current word
			w.addWord()
			_, _ = w.space.WriteRune(c)
		} else if inGroup(w.Breakpoints, c) {
			// valid breakpoint
			w.addSpace()
			w.addWord()
			_, _ = w.word.WriteRune(c)

			// Wrap line if the breakpoint would exceed the Limit
			if w.HardWrap && w.lineLen+w.space.Len()+runewidth.RuneWidth(c) > w.Limit {
				w.addNewLine()
			}

			// treat breakpoint as single character length words
			w.addWord()
		} else if w.HardWrap && w.lineLen+w.word.PrintableRuneWidth()+runewidth.RuneWidth(c)+w.space.Len() == w.Limit {
			// Word is at the limit -> begin new word
			_, _ = w.word.WriteRune(c)
			w.addWord()
		} else {
			// any other character
			_, _ = w.word.WriteRune(c)

			// add a line break if the current word would exceed the line's
			// character limit
			if w.lineLen+w.space.Len()+w.word.PrintableRuneWidth() > w.Limit &&
				w.word.PrintableRuneWidth() < w.Limit {
				w.addNewLine()
			}
		}
	}

	return len(b), nil
}

// Close will finish the word-wrap operation. Always call it before trying to
// retrieve the final result.
func (w *WordWrap) Close() error {
	if w.PreserveSpaces {
		w.addSpace()
	}
	w.addWord()

	return nil
}

// Bytes returns the word-wrapped result as a byte slice.
// Make sure to have closed the wordwrapper, before calling it.
func (w *WordWrap) Bytes() []byte {
	return w.buf.Bytes()
}

// String returns the word-wrapped result as a string.
// Make sure to have closed the wordwrapper, before calling it.
func (w *WordWrap) String() string {
	return w.buf.String()
}
