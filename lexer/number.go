package lexer

import "regexp"

var numberRegex = regexp.MustCompile(`\A([0-9]+)`)

// Matching numbers.
func findNumber(code string) Token {
	return find(constantRegex, code, "NUMBER")
}
