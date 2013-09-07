// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

var stringRegex = regexp.MustCompile(`\A"(.*?)"`)

// Matching strings.
func findString(code string) Token {
	token := findSubmatch(stringRegex, code, "STRING")
	token.length += 2
	return token
}
