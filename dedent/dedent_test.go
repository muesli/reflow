package dedent

import (
	"testing"
)

func TestDedent(t *testing.T) {
	tt := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "      --help      Show help for command\n      --version   Show version\n",
			Expected: "--help      Show help for command\n--version   Show version\n",
		},
		{
			Input:    "      --help              Show help for command\n  -C, --config string   Specify the config file to use\n",
			Expected: "    --help              Show help for command\n-C, --config string   Specify the config file to use\n",
		},
		{
			Input:    "  line 1\n\n  line 2\n line 3",
			Expected: " line 1\n\n line 2\nline 3",
		},
		{
			Input:    "  line 1\n  line 2\n  line 3\n\n",
			Expected: "line 1\nline 2\nline 3\n\n",
		},
		{
			Input:    " \tline 1\n\t\tline 2\n\t line 3\n\n",
			Expected: "line 1\nline 2\nline 3\n\n",
		},
		{
			Input:    "\t\tline 1\n\n\t\tline 2\n\tline 3",
			Expected: "\tline 1\n\n\tline 2\nline 3",
		},
		{
			Input:    "\n\n\n\n\n\n",
			Expected: "\n\n\n\n\n\n",
		},
		{
			Input:    "",
			Expected: "",
		},
		// Issue #72: no shared indent should not strip spaces within lines
		{
			Input:    "This is an example\n\t\n  with no shared indent.",
			Expected: "This is an example\n\t\n  with no shared indent.",
		},
		// Zero indent on first line, followed by indented lines
		{
			Input:    "First line\n  second line\n    third line",
			Expected: "First line\n  second line\n    third line",
		},
		// Indented first line, followed by line with zero indent
		{
			Input:    "  first line\nSecond line\n  third line",
			Expected: "  first line\nSecond line\n  third line",
		},
		// Varied indentation levels with common prefix
		{
			Input:    "    level 2\n      level 3\n    level 2 again\n        level 4",
			Expected: "level 2\n  level 3\nlevel 2 again\n    level 4",
		},
		// Internal spaces preserved across lines
		{
			Input:    "   foo   bar   baz\n   qux   quux",
			Expected: "foo   bar   baz\nqux   quux",
		},
		// Blank and empty lines between indented blocks
		{
			Input:    "   alpha\n\n     \n   beta\n \n   gamma",
			Expected: "alpha\n\n  \nbeta\n\ngamma",
		},
		// Blank lines with tabs
		{
			Input:    "  alpha\n\t\n  beta",
			Expected: "alpha\n\nbeta",
		},
		// Mixed tab and space indentation with preservation
		{
			Input:    "\t  alpha\n\t    beta\n\t  gamma",
			Expected: "alpha\n  beta\ngamma",
		},
		// Multiple tabs with relative indent preserved
		{
			Input:    "\t\talpha\n\t\t\tbeta\n\t\tgamma",
			Expected: "alpha\n\tbeta\ngamma",
		},
		// Windows CRLF line endings
		{
			Input:    "  hello\r\n\r\n    world\r\n  !\r\n",
			Expected: "hello\r\n\r\n  world\r\n!\r\n",
		},
		// Whitespace only lines
		{
			Input:    "   \n   \n",
			Expected: "   \n   \n",
		},
		// Single line with indentation
		{
			Input:    "   single line",
			Expected: "single line",
		},
		// Single line without indentation
		{
			Input:    "single line",
			Expected: "single line",
		},
		// No trailing newline
		{
			Input:    "  alpha\n    beta",
			Expected: "alpha\n  beta",
		},
	}

	for i, tc := range tt {
		s := String(tc.Input)
		if s != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Expected, s)
		}
	}
}

// go test -bench=BenchmarkDedent -benchmem -count=4
func BenchmarkDedent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		input := "  line 1\n\n  line 2\n line 3"
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			String(input)
		}
	})
}
