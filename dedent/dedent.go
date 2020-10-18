// +build !go1.10

package dedent

import (
	"bytes"
	"strings"
)

// String automatically detects the maximumen identation shared by all lines and
// trims them accordingly.
func String(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		indent := len(l) - len(strings.TrimLeft(l, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buf bytes.Buffer
	for _, l := range lines {
		l = strings.TrimPrefix(l, strings.Repeat(" ", minIndent))
		buf.WriteString(l + "\n")
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
