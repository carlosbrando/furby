// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

var constantRegex = regexp.MustCompile(`\A([A-Z]\w*)`)

// Matching class names and constants starting with a capital letter.
func findConstant(code string) Token {
	return find(constantRegex, code, "CONSTANT")
}
