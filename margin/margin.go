package margin

import (
	"io"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/padding"
)

type Writer struct {
	w  io.Writer
	pw *padding.Writer
	iw *indent.Writer
}

func NewWriter(w io.Writer, width uint, margin uint, marginFunc func(io.Writer)) *Writer {
	pw := &padding.Writer{
		Padding: width,
		PadFunc: marginFunc,
		Forward: &ansi.Writer{
			Forward: w,
		},
	}
	iw := &indent.Writer{
		Indent:     margin,
		IndentFunc: marginFunc,
		Forward: &ansi.Writer{
			Forward: pw,
		},
	}

	return &Writer{
		w:  w,
		pw: pw,
		iw: iw,
	}
}

func (w *Writer) Write(b []byte) (int, error) {
	return w.iw.Write(b)
}
