package lexer

import "regexp"

var constantRegex = regexp.MustCompile(`\A([A-Z]\w*)`)

// Matching class names and constants starting with a capital letter.
func findConstant(code string) Token {
	return find(constantRegex, code, "CONSTANT")
}
