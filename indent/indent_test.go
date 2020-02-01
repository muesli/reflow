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

func TestIndentString(t *testing.T) {
	actual := String("foobar", 3)
	expected := "   foobar"
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}
