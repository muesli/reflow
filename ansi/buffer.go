package ansi

import (
	"bytes"
	"unicode/utf8"

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
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	const (
		text int8 = iota + 1
		nF
		csi
		stTerminated
		// terminology
	)
	// var recognizeTerminologyEscSequences bool = false
	var state int8 = text
	var n int
	var cpIdx int = -1 // index of first code point byte
	// range over bytes not runes
	for i := 0; i < len(s); i++ {
		switch state {
		case nF: // [0x20—0x2F]+[0x30-0x7E]
			switch {
			case s[i] >= 0x20 && s[i] <= 0x2F:
			case i > 0 && s[i-1] >= 0x20 && s[i-1] <= 0x2F && s[i] >= 0x30 && s[i] <= 0x7E:
				state = text
			default:
				// fail
				state = text
			}
		case csi: // [0x40-0x7E]
			if s[i] >= 0x40 && s[i] <= 0x7E {
				state = text
			}
		case stTerminated:
			if s[i] == '\a' || (i > 0 && s[i] == '\\' && s[i-1] == 0x1B) {
				state = text
			}
		// case terminology:
		//     if s[i] == '\x00' {
		//         state = text
		//     }
		case text:
			if i > 0 && s[i-1] == 0x1B {
				switch {
				// nF escape sequences [0x20—0x2F]+[0x30-0x7E]
				case s[i] >= 0x20 && s[i] <= 0x2F:
					state = nF

				// Fp escape sequences [0x30—0x3F]
				case s[i] >= 0x30 && s[i] <= 0x3F:

				// Fe escape sequences [0x40-0x5F]
				// CSI - terminated by [0x40-0x7E]
				case s[i] == '[':
					state = csi
				// DCS, OSC, SOS, PM, APC - ST terminated
				case s[i] == 'P' || s[i] == ']' || s[i] == 'X' || s[i] == '^' || s[i] == '_':
					state = stTerminated
				case s[i] >= 0x40 && s[i] <= 0x5F:

				// Terminology  \x00 terminated - conflicts with Fs escape sequences
				// https://github.com/borisfaure/terminology/tree/master#extended-escapes-for-terminology-only
				// case recognizeTerminologyEscSequences && s[i] == '}':
				//     state = terminology
				// Fs escape sequences [0x60—0x7E]
				case s[i] >= 0x60 && s[i] <= 0x7E:
				}
			} else {
				if utf8.RuneStart(s[i]) {
					cpIdx = i
				}
				// cpBytes := unsafe.Slice(unsafe.StringData(s[cpIdx:i+1]), i-cpIdx+1) // go 1.20
				cpBytes := []byte(s[cpIdx : i+1])
				if utf8.FullRune(cpBytes) {
					if rn, _ := utf8.DecodeRune(cpBytes); rn != utf8.RuneError {
						n += runewidth.RuneWidth(rn)
					}
				}
			}
		}
	}

	return n
}
