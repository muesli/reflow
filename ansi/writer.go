package ansi

import (
	"bytes"
	"io"
)

type Writer struct {
	Forward io.Writer

	ansi       bool
	ansiseq    bytes.Buffer
	lastseq    bytes.Buffer
	seqchanged bool
}

// Write is used to write content to the ANSI buffer.
func (w *Writer) Write(b []byte) (int, error) {
	for i, c := range b {
		if c == '\x1B' {
			// ANSI escape sequence
			w.ansi = true
			w.seqchanged = true
			_ = w.ansiseq.WriteByte(c)
		} else if w.ansi {
			_ = w.ansiseq.WriteByte(c)
			if (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				w.ansi = false

				if bytes.HasSuffix(w.ansiseq.Bytes(), []byte("[0m")) {
					// reset sequence
					w.lastseq.Reset()
					w.seqchanged = false
				} else if c == 'm' {
					// color code
					_, _ = w.lastseq.Write(w.ansiseq.Bytes())
				}

				_, _ = w.ansiseq.WriteTo(w.Forward)
			}
		} else {
			_, err := w.Forward.Write(b[i : i+1])
			if err != nil {
				return 0, err
			}
		}
	}

	return len(b), nil
}

func (w *Writer) LastSequence() string {
	return w.lastseq.String()
}

func (w *Writer) ResetAnsi() {
	if !w.seqchanged {
		return
	}
	_, _ = w.Forward.Write([]byte("\x1b[0m"))
}

func (w *Writer) RestoreAnsi() {
	_, _ = w.Forward.Write(w.lastseq.Bytes())
}
