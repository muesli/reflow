package margin

import (
	"testing"
)

func TestMargin(t *testing.T) {
	tt := []struct {
		Input    string
		Expected string
		Width    uint
		Margin   uint
	}{
		// No-op, should pass through:
		{
			"foobar",
			"foobar",
			0,
			0,
		},
		// Basic margin:
		{
			"foobar",
			"  foobar  ",
			10,
			2,
		},
		// Multi-line margin:
		{
			"foo\nbar",
			" foo \n bar ",
			5,
			1,
		},
		// Don't pad empty trailing lines:
		{
			"foo\nbar\n",
			" foo \n bar \n",
			5,
			1,
		},
		// ANSI sequence codes:
		{
			"\x1B[38;2;249;38;114mfoo",
			"\x1B[38;2;249;38;114m\x1B[0m   \x1B[38;2;249;38;114mfoo   ",
			9,
			3,
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Width, tc.Margin, nil)

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

func TestMarginString(t *testing.T) {
	actual := String("foobar", 10, 2)
	expected := "  foobar  "
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}
