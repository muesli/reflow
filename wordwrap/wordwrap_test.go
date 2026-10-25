package wordwrap

import (
	"testing"

	"github.com/muesli/reflow/ansi"
)

func TestWordWrap(t *testing.T) {
	tt := []struct {
		Input        string
		Expected     string
		Limit        int
		KeepNewlines bool
	}{
		// No-op, should pass through, including trailing whitespace:
		{
			"foobar\n ",
			"foobar\n ",
			0,
			true,
		},
		// Nothing to wrap here, should pass through:
		{
			"foo",
			"foo",
			4,
			true,
		},
		// A single word that is too long passes through.
		// We do not break long words:
		{
			"foobarfoo",
			"foobarfoo",
			4,
			true,
		},
		// Lines are broken at whitespace:
		{
			"foo bar foo",
			"foo\nbar\nfoo",
			4,
			true,
		},
		// A hyphen is a valid breakpoint:
		{
			"foo-foobar",
			"foo-\nfoobar",
			4,
			true,
		},
		// Space buffer needs to be emptied before breakpoints:
		{
			"foo --bar",
			"foo --bar",
			9,
			true,
		},
		// Lines are broken at whitespace, even if words
		// are too long. We do not break words:
		{
			"foo bars foobars",
			"foo\nbars\nfoobars",
			4,
			true,
		},
		// A word that would run beyond the limit is wrapped:
		{
			"foo bar",
			"foo\nbar",
			5,
			true,
		},
		// Whitespace that trails a line and fits the width
		// passes through, as does whitespace prefixing an
		// explicit line break. A tab counts as one character:
		{
			"foo\nb\t a\n bar",
			"foo\nb\t a\n bar",
			4,
			true,
		},
		// Trailing whitespace is removed if it doesn't fit the width.
		// Runs of whitespace on which a line is broken are removed:
		{
			"foo    \nb   ar   ",
			"foo\nb\nar",
			4,
			true,
		},
		// An explicit line break at the end of the input is preserved:
		{
			"foo bar foo\n",
			"foo\nbar\nfoo\n",
			4,
			true,
		},
		// Explicit break are always preserved:
		{
			"\nfoo bar\n\n\nfoo\n",
			"\nfoo\nbar\n\n\nfoo\n",
			4,
			true,
		},
		// Unless we ask them to be ignored:
		{
			"\nfoo bar\n\n\nfoo\n",
			"foo\nbar\nfoo",
			4,
			false,
		},
		// Complete example:
		{
			" This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* foo  \nbar    ",
			" This\nis a\nlist: \n\n\t* foo\n\t* bar\n\n\n\t* foo\nbar",
			6,
			true,
		},
		// ANSI sequence codes don't affect length calculation:
		{
			"\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m",
			"\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m",
			7,
			true,
		},
		// ANSI control codes don't get wrapped:
		{
			"\x1B[38;2;249;38;114m(\x1B[0m\x1B[38;2;248;248;242mjust another test\x1B[38;2;249;38;114m)\x1B[0m",
			"\x1B[38;2;249;38;114m(\x1B[0m\x1B[38;2;248;248;242mjust\nanother\ntest\x1B[38;2;249;38;114m)\x1B[0m",
			3,
			true,
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Limit)
		f.KeepNewlines = tc.KeepNewlines

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

func TestWordWrapString(t *testing.T) {
	actual := String("foo bar", 3)
	expected := "foo\nbar"
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}

func TestWordWrapUnicodeWhitespaceWidths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		limit    int
		width    int
		expected string
	}{
		{"em_space_fits", "a\u2003b", 3, 3, "a\u2003b"},
		{"em_space_wraps", "a\u2003b", 2, 3, "a\nb"},
		{"repeated_em_spaces_fit", "a\u2003\u2003b", 4, 4, "a\u2003\u2003b"},
		{"repeated_em_spaces_wrap", "a\u2003\u2003b", 3, 4, "a\nb"},
		{"ideographic_space_fits", "a\u3000b", 4, 4, "a\u3000b"},
		{"ideographic_space_wraps", "a\u3000b", 3, 4, "a\nb"},
		{"mixed_spaces_fit", "a\u2003\u3000b", 5, 5, "a\u2003\u3000b"},
		{"mixed_spaces_wrap", "a\u2003\u3000b", 4, 5, "a\nb"},
		{"emitted_space_width", "a\u2003b c", 5, 5, "a\u2003b c"},
		{"nbsp_fits", "a\u00a0b", 3, 3, "a\u00a0b"},
		{"narrow_nbsp_fits", "a\u202fb", 3, 3, "a\u202fb"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if width := ansi.PrintableRuneWidth(tc.input); width != tc.width {
				t.Fatalf("fixture width = %d, want %d", width, tc.width)
			}
			if actual := String(tc.input, tc.limit); actual != tc.expected {
				t.Errorf("String(%q, %d) = %q, want %q", tc.input, tc.limit, actual, tc.expected)
			}
			if actual := string(Bytes([]byte(tc.input), tc.limit)); actual != tc.expected {
				t.Errorf("Bytes(%q, %d) = %q, want %q", tc.input, tc.limit, actual, tc.expected)
			}
		})
	}
}

func TestWordWrapWhitespaceCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		limit    int
		expected string
	}{
		{"space_fits", "a b", 3, "a b"},
		{"space_wraps", "a b", 2, "a\nb"},
		{"tab_fits", "a\tb", 3, "a\tb"},
		{"tab_wraps", "a\tb", 2, "a\nb"},
		{"carriage_return_fits", "a\rb", 3, "a\rb"},
		{"carriage_return_wraps", "a\rb", 2, "a\nb"},
		{"vertical_tab_fits", "a\vb", 3, "a\vb"},
		{"vertical_tab_wraps", "a\vb", 2, "a\nb"},
		{"form_feed_fits", "a\fb", 3, "a\fb"},
		{"form_feed_wraps", "a\fb", 2, "a\nb"},
		{"next_line_fits", "a\u0085b", 4, "a\u0085b"},
		{"next_line_wraps", "a\u0085b", 3, "a\nb"},
		{"nbsp_remains_breakable", "a\u00a0b", 2, "a\nb"},
		{"narrow_nbsp_remains_breakable", "a\u202fb", 2, "a\nb"},
		{"limit_zero", "\x1b[31ma\u2003b\x1b[0m\n\t ", 0, "\x1b[31ma\u2003b\x1b[0m\n\t "},
		{"long_word", "foobarfoo", 4, "foobarfoo"},
		{"long_word_after_space", "a foobarfoo", 4, "a\nfoobarfoo"},
		{"limit_one", "a b", 1, "a b"},
		{"trailing_whitespace", "a \u2003\t", 8, "a"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := String(tc.input, tc.limit); actual != tc.expected {
				t.Errorf("String(%q, %d) = %q, want %q", tc.input, tc.limit, actual, tc.expected)
			}
		})
	}
}

func TestWordWrapUnicodeWhitespaceANSI(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected string
	}{
		{"fits", 3, "\x1b[31ma\u2003b\x1b[0m"},
		{"wraps", 2, "\x1b[31ma\nb\x1b[0m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if actual := String("\x1b[31ma\u2003b\x1b[0m", tc.limit); actual != tc.expected {
				t.Errorf("ANSI output = %q, want %q", actual, tc.expected)
			}
		})
	}
}

func TestWordWrapUnicodeWhitespaceNewlines(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		limit        int
		newline      []rune
		keepNewlines bool
		beforeClose  string
		expected     string
	}{
		{"em_space_fits", "a\u2003\nb", 2, nil, true, "a\u2003\n", "a\u2003\nb"},
		{"ideographic_space_fits", "a\u3000\nb", 3, nil, true, "a\u3000\n", "a\u3000\nb"},
		{"ideographic_space_dropped", "a\u3000\nb", 2, nil, true, "a\n", "a\nb"},
		{"line_state_reset", "a\u2003\nb\u2003c", 3, nil, true, "a\u2003\nb", "a\u2003\nb\u2003c"},
		{"custom_newline", "a\u2003\rb", 2, []rune{'\r'}, true, "a\u2003\n", "a\u2003\nb"},
		{"ascii_newline", "a \nb", 2, nil, true, "a \n", "a \nb"},
		{"ignore_newline", "a\nb", 3, nil, false, "a", "a b"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWriter(tc.limit)
			w.KeepNewlines = tc.keepNewlines
			if tc.newline != nil {
				w.Newline = tc.newline
			}
			assertWordWrapWrite(t, w, tc.input, tc.beforeClose)
			assertWordWrapClose(t, w, tc.expected)
		})
	}
}

func TestWordWrapWhitespaceBreakpoints(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		limit       int
		breakpoints []rune
		expected    string
	}{
		{"unicode_space", "a\u2003:b", 4, []rune{':'}, "a\u2003:b"},
		{"historical_breakpoint_width", "a:b", 2, []rune{':'}, "a:b"},
		{"default_hyphen", "foo-foobar", 4, nil, "foo-\nfoobar"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWriter(tc.limit)
			if tc.breakpoints != nil {
				w.Breakpoints = tc.breakpoints
			}
			if _, err := w.Write([]byte(tc.input)); err != nil {
				t.Fatal(err)
			}
			assertWordWrapClose(t, w, tc.expected)
		})
	}
}

func TestWordWrapWhitespaceStreaming(t *testing.T) {
	tests := []struct {
		name         string
		limit        int
		keepNewlines bool
		chunks       []string
		intermediate []string
		expected     string
	}{
		{"unicode_fits", 3, true, []string{"a", "\u2003", "b"}, []string{"", "a", "a"}, "a\u2003b"},
		{"unicode_wraps", 2, true, []string{"a", "\u2003", "b"}, []string{"", "a", "a\n"}, "a\nb"},
		{"split_ansi", 3, true, []string{"\x1b[", "31ma", "\u2003", "b\x1b[", "0m"}, []string{"", "", "\x1b[31ma", "\x1b[31ma", "\x1b[31ma"}, "\x1b[31ma\u2003b\x1b[0m"},
		{"after_automatic_wrap", 2, true, []string{"a", "\u2003", "b", " c"}, []string{"", "a", "a\n", "a\nb\n"}, "a\nb\nc"},
		{"empty_writes", 3, true, []string{"", "a", "", "\u2003", "", "b", ""}, []string{"", "", "", "a", "a", "a", "a"}, "a\u2003b"},
		{"per_write_trimming", 3, false, []string{"a ", " b"}, []string{"", ""}, "ab"},
		{"split_utf8", 3, true, []string{"a\xe2", "\x80\x83b"}, []string{"", ""}, "a\ufffd\ufffd\ufffdb"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := NewWriter(tc.limit)
			w.KeepNewlines = tc.keepNewlines
			for i, chunk := range tc.chunks {
				assertWordWrapWrite(t, w, chunk, tc.intermediate[i])
			}
			assertWordWrapClose(t, w, tc.expected)
			assertWordWrapClose(t, w, tc.expected)
		})
	}
}

func assertWordWrapWrite(t *testing.T, w *WordWrap, input, expected string) {
	t.Helper()
	n, err := w.Write([]byte(input))
	if err != nil {
		t.Errorf("Write(%q) error = %v", input, err)
	}
	if n != len(input) {
		t.Errorf("Write(%q) count = %d, want %d", input, n, len(input))
	}
	if actual := w.String(); actual != expected {
		t.Errorf("after Write(%q) = %q, want %q", input, actual, expected)
	}
}

func assertWordWrapClose(t *testing.T, w *WordWrap, expected string) {
	t.Helper()
	if err := w.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	if actual := w.String(); actual != expected {
		t.Errorf("after Close() = %q, want %q", actual, expected)
	}
}
