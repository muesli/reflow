package truncate

import (
	"bytes"
	"errors"
	"testing"

	"github.com/muesli/reflow/ansi"
)

func TestTruncate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		width    uint
		tail     string
		in       string
		expected string
	}{
		// No-op, should pass through:
		{
			10,
			"",
			"hello",
			"hello",
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.width, tc.tail)

		_, err := f.Write([]byte(tc.in))
		if err != nil {
			t.Error(err)
		}

		if f.String() != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.expected, f.String())
		}
	}
}

func TestTruncateString(t *testing.T) {
	t.Parallel()

	actual := String("foobar", 3)
	expected := "foo"
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func BenchmarkTruncateString(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			String("foo", 2)
		}
	})
}

func TestTruncateWriterWithTail(t *testing.T) {
	t.Parallel()

	f := NewWriter(5, "...")

	_, err := f.Write([]byte("foo\n"))
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write([]byte("bar"))
	if err != nil {
		t.Error(err)
	}

	exp := "..foo\n..bar"
	if f.String() != exp {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", exp, f.String())
	}
}

func TestNewWriterPipe(t *testing.T) {
	t.Parallel()

	b := &bytes.Buffer{}
	f := NewWriterPipe(b, 2, "")

	if _, err := f.Write([]byte("foo")); err != nil {
		t.Error(err)
	}

	actual := b.String()
	expected := "fo"

	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func TestWriter_Error(t *testing.T) {
	t.Parallel()

	f := &Writer{
		ansiWriter: &ansi.Writer{Forward: fakeWriter{}},
	}

	if _, err := f.Write([]byte("foo")); err != fakeErr {
		t.Error(err)
	}

	if _, err := f.Write([]byte("foo")); err != fakeErr {
		t.Error(err)
	}
}

var fakeErr = errors.New("fake error")

type fakeWriter struct{}

func (fakeWriter) Write(_ []byte) (int, error) {
	return 0, fakeErr
}
