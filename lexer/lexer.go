// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"regexp"
	"strings"
)

type Token struct {
	tokenType string
	value     string
	length    int
	matched   bool
}

var whitespaceRegex = regexp.MustCompile(`\A\s+`)

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

		// ...

		if !token.matched {
			token = findOperator(chunk)
		}

		if !token.matched {
			if m := whitespaceRegex.FindString(chunk); m != "" {
				i += len(m)
				continue // ignore whitespace
			}
		}

		// catch all single characters
		// we treat all other single characters as a token. Eg.: ( ) , . ! + - <
		if !token.matched {
			value := chunk[:1]
			token = Token{value, value, 1, true}
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
