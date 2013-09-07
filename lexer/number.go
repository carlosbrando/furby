// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

var numberRegex = regexp.MustCompile(`\A([0-9]+)`)

// Matching numbers.
func findNumber(code string) Token {
	return find(numberRegex, code, "NUMBER")
}
