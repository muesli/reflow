package padding

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/muesli/reflow/ansi"
)

func TestPadding(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Input    string
		Expected string
		Padding  uint
	}{
		// No-op, should pass through:
		{
			"foobar",
			"foobar",
			0,
		},
		// Basic padding:
		{
			"foobar",
			"foobar    ",
			10,
		},
		// Multi-line padding:
		{
			"foo\nbar",
			"foo   \nbar   ",
			6,
		},
		// Don't pad empty trailing lines:
		{
			"foo\nbar\n",
			"foo   \nbar   \n",
			6,
		},
		// ANSI sequence codes:
		{
			"\x1B[38;2;249;38;114mfoo",
			"\x1B[38;2;249;38;114mfoo   ",
			6,
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Padding, nil)

		_, err := f.Write([]byte(tc.Input))
		if err != nil {
			t.Error(err)
		}

		if err := f.Close(); err != nil {
			t.Error(err)
		}

		if f.String() != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Expected, f.String())
		}
	}
}

func TestPaddingWriter(t *testing.T) {
	t.Parallel()

	f := NewWriter(6, nil)

	_, err := f.Write([]byte("foo\n"))
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write([]byte("bar"))
	if err != nil {
		t.Error(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}

	exp := "foo   \nbar   "
	if f.String() != exp {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", exp, f.String())
	}
}

func TestPaddingString(t *testing.T) {
	t.Parallel()

	actual := String("foobar", 10)
	expected := "foobar    "
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func BenchmarkPaddingString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			String("foobar", 10)
		}
	})
}

func TestNewWriterPipe(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}
	f := NewWriterPipe(b, 10, nil)

	if _, err := f.Write([]byte("foobar")); err != nil {
		t.Error(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}

	actual := b.String()
	expected := "foobar    "

	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func TestWriter_pad(t *testing.T) {
	t.Parallel()

	f := NewWriter(4, func(w io.Writer) {
		_, _ = w.Write([]byte("."))
	})

	if err := f.pad(); err != nil {
		t.Error(err)
	}

	actual := f.buf.String()
	expected := "...."
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func TestWriter_Flush(t *testing.T) {
	t.Parallel()

	f := NewWriter(6, nil)

	_, err := f.Write([]byte("foo"))
	if err != nil {
		t.Error(err)
	}

	if err := f.Flush(); err != nil {
		t.Error(err)
	}

	exp := "foo   "
	if f.String() != exp {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", exp, f.String())
	}

	_, err = f.Write([]byte("bar"))
	if err != nil {
		t.Error(err)
	}
	if err := f.Flush(); err != nil {
		t.Error(err)
	}

	exp = "bar   "
	if f.String() != exp {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", exp, f.String())
	}
}

func TestWriter_Close(t *testing.T) {
	t.Parallel()

	f := &Writer{
		Padding:    6,
		lineLen:    1,
		ansiWriter: &ansi.Writer{Forward: fakeWriter{}},
	}

	if err := f.Close(); err != fakeErr {
		t.Error(err)
	}
}

func TestWriter_Error(t *testing.T) {
	t.Parallel()

	f := &Writer{
		Padding:    6,
		ansiWriter: &ansi.Writer{Forward: fakeWriter{}},
	}

	if _, err := f.Write([]byte("foo\n")); err != fakeErr {
		t.Error(err)
	}

	if _, err := f.Write([]byte("\n")); err != fakeErr {
		t.Error(err)
	}

	if err := f.pad(); err != fakeErr {
		t.Error(err)
	}
}

var fakeErr = errors.New("fake error")

type fakeWriter struct{}

func (fakeWriter) Write(_ []byte) (int, error) {
	return 0, fakeErr
}

func TestPaddingLargeWidth(t *testing.T) {
	t.Parallel()

	largeWidths := []uint{
		^uint(0),
		^uint(0) - 1,
		uint(maxInt) + 1,
	}

	for _, w := range largeWidths {
		res := String("hello", w)
		expected := "hello"
		if res != expected {
			t.Errorf("expected %q, got %q for width %d", expected, res, w)
		}

		bRes := Bytes([]byte("hello"), w)
		if string(bRes) != expected {
			t.Errorf("expected %q, got %q for width %d", expected, string(bRes), w)
		}

		var buf bytes.Buffer
		pipe := NewWriterPipe(&buf, w, nil)
		if _, err := pipe.Write([]byte("hello\nworld\n")); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := pipe.Close(); err != nil {
			t.Errorf("unexpected error on close: %v", err)
		}
		if buf.String() != "hello\nworld\n" {
			t.Errorf("expected %q, got %q for width %d", "hello\nworld\n", buf.String(), w)
		}

		var called bool
		writerWithFunc := NewWriter(w, func(wr io.Writer) {
			called = true
			_, _ = wr.Write([]byte(">"))
		})
		if _, err := writerWithFunc.Write([]byte("hello")); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if err := writerWithFunc.Close(); err != nil {
			t.Errorf("unexpected error on close: %v", err)
		}
		if called {
			t.Errorf("pad func should not be called for overflowed width %d", w)
		}
		if writerWithFunc.String() != expected {
			t.Errorf("expected %q, got %q for width %d", expected, writerWithFunc.String(), w)
		}
	}
}
