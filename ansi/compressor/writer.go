package compressor

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/muesli/reflow/ansi"
)

type Writer struct {
	Forward io.Writer

	ansi        bool
	ansiseq     bytes.Buffer
	lastseq     bytes.Buffer
	prevlastseq bytes.Buffer
	resetreq    bool
	runeBuf     []byte
}

// Write is used to write content to the ANSI buffer.
func (w *Writer) Write(b []byte) (int, error) {
	for _, c := range string(b) {
		if c == ansi.Marker {
			// ANSI escape sequence
			w.ansi = true
			_, _ = w.ansiseq.WriteRune(c)
		} else if w.ansi {
			_, _ = w.ansiseq.WriteRune(c)
			if ansi.IsTerminator(c) {
				// ANSI sequence terminated
				w.ansi = false

				terminated := false
				if bytes.HasSuffix(w.ansiseq.Bytes(), []byte("[0m")) {
					// reset sequence
					w.prevlastseq.Reset()
					w.prevlastseq.Write(w.lastseq.Bytes())

					w.lastseq.Reset()
					terminated = true
					w.resetreq = true
				} else if c == 'm' {
					// color code
					_, _ = w.lastseq.Write(w.ansiseq.Bytes())
				}

				if !terminated {
					// did we reset the sequence just to restore it again?
					if bytes.Equal(w.ansiseq.Bytes(), w.prevlastseq.Bytes()) {
						w.resetreq = false
						w.ansiseq.Reset()
					}

					w.prevlastseq.Reset()

					if w.resetreq {
						w.ResetAnsi()
					}

					_, _ = w.Forward.Write(w.ansiseq.Bytes())
				}

				w.ansiseq.Reset()
			}
		} else {
			if w.resetreq {
				w.ResetAnsi()
			}

			_, err := w.writeRune(c)
			if err != nil {
				return 0, err
			}
		}
	}

	return len(b), nil
}

func (w *Writer) writeRune(r rune) (int, error) {
	if w.runeBuf == nil {
		w.runeBuf = make([]byte, utf8.UTFMax)
	}
	n := utf8.EncodeRune(w.runeBuf, r)
	return w.Forward.Write(w.runeBuf[:n])
}

// Close finishes the compression operation. Always call it before trying to
// retrieve the final result.
func (w *Writer) Close() error {
	if w.resetreq {
		w.ResetAnsi()
	}

	return nil
}

func (w *Writer) ResetAnsi() {
	w.prevlastseq.Reset()
	_, _ = w.Forward.Write([]byte("\x1b[0m"))
	w.resetreq = false
}
