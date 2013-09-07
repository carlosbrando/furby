// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

var identifierRegex = regexp.MustCompile(`\A([a-z]\w*)`)

var keywords = map[string]string{
	"def":   "DEF",
	"if":    "IF",
	"true":  "TRUE",
	"false": "FALSE",
	"nil":   "NIL",
}

// Matching if, print, method names, etc.
// Keywords are special identifiers tagged with their own name, 'if' will result in an [:IF, "if"] token.
// Non-keyword identifiers include method and variable names.
func findIdentifier(code string) Token {
	token := find(identifierRegex, code, "IDENTIFIER")

	if tokenType, ok := keywords[token.value]; ok {
		token.tokenType = tokenType
	}

	return token
}
