package ansi

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	buf := []byte("\x1B[38;2;249;38;114mfoo\x1B[0m")
	w := &Writer{Forward: ioutil.Discard}

	n, err := w.Write(buf)

	w.ResetAnsi()

	w.RestoreAnsi()

	if err != nil {
		t.Fatalf("err should be nil, but got %v", err)
	}

	if l := len(buf); n != l {
		t.Fatalf("n should be %d, got %d", l, n)
	}

	if ls := w.LastSequence(); ls != "" {
		t.Fatalf("LastSequence should be empty, got %s", ls)
	}
}

var fakeErr = errors.New("fake error")

type fakeWriter struct{}

func (fakeWriter) Write(_ []byte) (int, error) {
	return 0, fakeErr
}

func TestWriter_Write_Error(t *testing.T) {
	w := &Writer{Forward: fakeWriter{}}

	_, err := w.Write([]byte("foo"))

	if err != fakeErr {
		t.Fatalf("err should be fakeErr, but got %v", err)
	}
}

func BenchmarkWriter_Write(b *testing.B) {
	buf := []byte("\x1B[38;2;249;38;114mfoo\x1B[0m")
	w := &Writer{Forward: ioutil.Discard}
	var (
		n   int
		err error
	)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err = w.Write(buf)
	}

	if err != nil {
		b.Fatalf("err should be nil, but got %v", err)
	}

	if l := len(buf); n != l {
		b.Fatalf("n should be %d, got %d", l, n)
	}

	if ls := w.LastSequence(); ls != "" {
		b.Fatalf("LastSequence should be empty, got %s", ls)
	}
}

func TestWriter_ResetAnsi(t *testing.T) {
	b := &bytes.Buffer{}
	w := &Writer{Forward: b}

	w.ResetAnsi()

	if b.String() != "" {
		t.Fatal("b should be empty")
	}

	w.seqchanged = true

	w.ResetAnsi()

	if s := b.String(); s != "\x1b[0m" {
		t.Fatalf("b.String() should be \"\\x1b[0m\", got %s", s)
	}
}

func TestWriter_RestoreAnsi(t *testing.T) {
	b := &bytes.Buffer{}
	w := &Writer{Forward: b, lastseq: "\x1B[38;2;249;38;114m"}

	w.RestoreAnsi()

	if s := b.String(); s != "\x1B[38;2;249;38;114m" {
		t.Fatalf("b.String() should be \"\\x1B[38;2;249;38;114m\", got %s", s)
	}
}
