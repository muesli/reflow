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
		f.Close()

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
	f.Close()

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
	var actual string
	expected := "foobar    "

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			actual = String("foobar", 10)
		}
	})

	if actual != expected {
		b.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func TestNewWriterPipe(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}
	f := NewWriterPipe(b, 10, nil)

	if _, err := f.Write([]byte("foobar")); err != nil {
		t.Error(err)
	}
	f.Close()

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
