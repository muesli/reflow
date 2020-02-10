package indent

import (
	"testing"
)

func TestIndent(t *testing.T) {
	tt := []struct {
		Input    string
		Expected string
		Indent   uint
	}{
		// No-op, should pass through:
		{
			"foobar",
			"foobar",
			0,
		},
		// Basic indentation:
		{
			"foobar",
			"    foobar",
			4,
		},
		// Multi-line indentation:
		{
			"foo\nbar",
			"    foo\n    bar",
			4,
		},
		// ANSI sequence codes:
		{
			"\x1B[38;2;249;38;114mfoo",
			"\x1B[38;2;249;38;114m\x1B[0m    \x1B[38;2;249;38;114mfoo",
			4,
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Indent, nil)

		_, err := f.Write([]byte(tc.Input))
		if err != nil {
			t.Error(err)
		}

		if f.String() != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Expected, f.String())
		}
	}
}

func TestIndentWriter(t *testing.T) {
	f := NewWriter(4, nil)

	_, err := f.Write([]byte("foo\n"))
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write([]byte("bar"))
	if err != nil {
		t.Error(err)
	}

	exp := "    foo\n    bar"
	if f.String() != exp {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", exp, f.String())
	}
}

func TestIndentString(t *testing.T) {
	actual := String("foobar", 3)
	expected := "   foobar"
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}
