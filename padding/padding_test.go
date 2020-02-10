package padding

import (
	"testing"
)

func TestPadding(t *testing.T) {
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
	actual := String("foobar", 10)
	expected := "foobar    "
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}
