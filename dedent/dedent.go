package dedent

import (
	"bytes"
)

// String automatically detects the maximum indentation shared by all lines and
// trims them accordingly.
func String(s string) string {
	indent := minIndent(s)
	if indent == 0 {
		return s
	}

	return dedent(s, indent)
}

func minIndent(s string) int {
	const unset = -1
	var (
		curIndent    int
		minIndent    = unset
		shouldAppend = true
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t':
			if shouldAppend {
				curIndent++
			}
		case '\n':
			curIndent = 0
			shouldAppend = true
		default:
			// First non-whitespace on this line. Update the minimum
			// even when curIndent is 0 so a flush-left line correctly
			// pulls the shared indent down to zero. Previously a zero
			// indent was indistinguishable from "no line seen yet"
			// and dedent silently kept stripping spaces off later
			// indented lines. shouldAppend is still true here, which
			// is how we know we're at the start of a content line.
			if shouldAppend && (minIndent == unset || curIndent < minIndent) {
				minIndent = curIndent
			}
			shouldAppend = false
		}
	}

	if minIndent < 0 {
		return 0
	}
	return minIndent
}

func dedent(s string, indent int) string {
	var (
		omitted    int
		shouldOmit = true
		buf        bytes.Buffer
	)

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t':
			if shouldOmit {
				if omitted < indent {
					omitted++
					continue
				}
				shouldOmit = false
			}
			_ = buf.WriteByte(s[i])
		case '\n':
			omitted = 0
			shouldOmit = true
			_ = buf.WriteByte(s[i])
		default:
			_ = buf.WriteByte(s[i])
		}
	}

	return buf.String()
}
