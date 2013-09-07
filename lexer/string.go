package lexer

import "regexp"

var stringRegex = regexp.MustCompile(`\A"(.*?)"`)

// Matching strings.
func findString(code string) Token {
	if m := stringRegex.FindStringSubmatch(code); len(m) > 0 {
		return Token{"STRING", m[1], len(m[1]) + 2, true}
	}

	return Token{"", "", 0, false}
}
