package lexer

import (
	"strings"
)

type Token struct {
	tokenType string
	value     string
	length    int
	matched   bool
}

func Tokenize(code string) [][]string {
	// cleanup code by remove extra line breaks
	code = strings.TrimSpace(code)

	// collection of all parsed tokens in the form [:TOKEN_TYPE, value]
	var tokens [][]string

	// We keep track of the indentation levels we are in so that when we dedent, we can
	// check if we're on the correct level.
	// indent_stack := []string

	// scan one character at the time until you find something to parse.
	for i := 0; i < len(code); {
		chunk := code[i:]

		// Matching standard tokens.
		//
		token := findIdentifier(chunk)

		if !token.matched {
			token = findConstant(chunk)
		}

		if !token.matched {
			token = findNumber(chunk)
		}

		if !token.matched {
			token = findString(chunk)
		}

		// if a token was found add it to the stack
		if token.matched {
			tokens = append(tokens, []string{token.tokenType, token.value})
			i += token.length
		} else {
			i += 1
		}
	}

	return tokens
}
