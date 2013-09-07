// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

var operatorRegex = regexp.MustCompile(`\A(\|\||&&|==|!=|<=|>=)`)

// Match long operators such as ||, &&, ==, !=, <= and >=.
// One character long operators are matched by the catch all regex.
func findOperator(code string) Token {
	token := find(operatorRegex, code, "OPERATOR")
	token.tokenType = token.value
	return token
}
