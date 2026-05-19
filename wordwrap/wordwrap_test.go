package wordwrap

import (
	"testing"
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

func TestWordWrapCJK(t *testing.T) {
	tt := []struct {
		Input    string
		Expected string
		Limit    int
	}{
		// Pure CJK: each character is a valid break point (2 cols each).
		// "中文"=4 cols, "测试"=4 cols → breaks at col 4.
		{
			"中文测试",
			"中文\n测试",
			4,
		},
		// "中文测"=6 cols, "试"=2 cols → breaks at col 6.
		{
			"中文测试",
			"中文测\n试",
			6,
		},
		// CJK fits in limit → no wrap.
		{
			"中文测试",
			"中文测试",
			8,
		},
		// CJK after Latin (space-separated): unchanged behavior.
		// "foo " = 4 cols, "中文" = 4 cols → 8 cols total, fits.
		{
			"foo 中文测试",
			"foo 中文\n测试",
			8,
		},
		// CJK attached to Latin (no space): break at CJK↔Latin boundary.
		// "这是"=4 cols, then "manual触" at newline → "manual触"=8 cols, fits 12.
		{
			"这是manual触发",
			"这是\nmanual触\n发",
			8,
		},
		// CJK punctuation: fullwidth punctuation is also a break point.
		{
			"你好，世界！",
			"你好，\n世界！",
			6,
		},
		// Mixed CJK+Latin with fullwidth punctuation.
		{
			"manual（手动触发），很可能没跑。",
			"manual（手动\n触发），很可\n能没跑。",
			12,
		},
		// Long CJK sentence wraps correctly.
		{
			"这是一个比较长的中文句子用来测试换行功能",
			"这是一个比\n较长的中文\n句子用来测\n试换行功能",
			10,
		},
		// Latin-only behavior unchanged.
		{
			"hello world",
			"hello\nworld",
			6,
		},
		// CJK↔Latin boundary: "测试"=4, then "abc"=3, then "测"=2 → "abc测"=5 > 4? No, limit=6.
		// "测试"=4 cols, then "abc" starts a new word (CJK↔Latin boundary).
		// "abc测" = 3+2=5 cols. "试" starts another word. 5+2=7 > 6 → break.
		{
			"测试abc测试",
			"测试\nabc测\n试",
			6,
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Limit)
		f.KeepNewlines = true

		_, err := f.Write([]byte(tc.Input))
		if err != nil {
			t.Error(err)
		}
		f.Close()

		if f.String() != tc.Expected {
			t.Errorf("Test %d (limit=%d, input=%q):\nexpected:\n\n`%s`\n\nActual Output:\n\n`%s`",
				i, tc.Limit, tc.Input, tc.Expected, f.String())
		}
	}
}

func TestWordWrapCJKNoWrap(t *testing.T) {
	// Limit 0 disables wrapping — CJK text passes through unchanged.
	got := String("中文测试abc测试", 0)
	want := "中文测试abc测试"
	if got != want {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", want, got)
	}
}

func TestWordWrapCJKString(t *testing.T) {
	got := String("你好世界", 4)
	want := "你好\n世界"
	if got != want {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", want, got)
	}
}
