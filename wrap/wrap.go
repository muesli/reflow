package wrap

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/muesli/reflow/ansi"
	"github.com/rivo/uniseg"
)

var (
	defaultNewline  = []rune{'\n'}
	defaultTabWidth = 4
)

type Wrap struct {
	Limit         int
	Newline       []rune
	KeepNewlines  bool
	PreserveSpace bool
	TabWidth      int

	buf             *bytes.Buffer
	lineLen         int
	ansi            bool
	forcefulNewline bool
}

// NewWriter returns a new instance of a wrapping writer, initialized with
// default settings.
func NewWriter(limit int) *Wrap {
	return &Wrap{
		Limit:        limit,
		Newline:      defaultNewline,
		KeepNewlines: true,
		// Keep whitespaces following a forceful line break. If disabled,
		// leading whitespaces in a line are only kept if the line break
		// was not forceful, meaning a line break that was already present
		// in the input
		PreserveSpace: false,
		TabWidth:      defaultTabWidth,

		buf: &bytes.Buffer{},
	}
}

// Bytes is shorthand for declaring a new default Wrap instance,
// used to immediately wrap a byte slice.
func Bytes(b []byte, limit int) []byte {
	f := NewWriter(limit)
	_, _ = f.Write(b)

	return f.Bytes()
}

func (w *Wrap) addNewLine() {
	_, _ = w.buf.WriteRune('\n')
	w.lineLen = 0
}

// String is shorthand for declaring a new default Wrap instance,
// used to immediately wrap a string.
func String(s string, limit int) string {
	return string(Bytes([]byte(s), limit))
}

func (w *Wrap) Write(b []byte) (int, error) {
	s := strings.Replace(string(b), "\t", strings.Repeat(" ", w.TabWidth), -1)
	if !w.KeepNewlines {
		s = strings.Replace(s, "\n", "", -1)
	}

	width := ansi.PrintableRuneWidth(s)

	if w.Limit <= 0 || w.lineLen+width <= w.Limit {
		w.lineLen += width
		return w.buf.Write(b)
	}

	state := -1
	var cluster string

	for len(s) > 0 {
		cluster, s, width, state = uniseg.FirstGraphemeClusterInString(s, state)
		rs := []rune(cluster)

		switch {
		case len(rs) == 1 && rs[0] == ansi.Marker:
			w.ansi = true
		case len(rs) == 1 && w.ansi && ansi.IsTerminator(rs[0]):
			w.ansi = false
		case w.ansi:
		case len(rs) == 1 && inGroup(w.Newline, rs[0]):
			w.addNewLine()
			w.forcefulNewline = false
			continue
		default:
			if w.lineLen+width > w.Limit {
				w.addNewLine()
				w.forcefulNewline = true
			}

			switch {
			case w.lineLen == 0:
				if len(rs) == 1 && w.forcefulNewline && !w.PreserveSpace && unicode.IsSpace(rs[0]) {
					continue
				}
			default:
				w.forcefulNewline = false
			}

			w.lineLen += width
		}

		_, _ = w.buf.WriteString(cluster)
	}

	return len(b), nil
}

// Bytes returns the wrapped result as a byte slice.
func (w *Wrap) Bytes() []byte {
	return w.buf.Bytes()
}

// String returns the wrapped result as a string.
func (w *Wrap) String() string {
	return w.buf.String()
}

func inGroup(a []rune, c rune) bool {
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}
