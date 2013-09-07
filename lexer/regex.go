// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import "regexp"

func find(r *regexp.Regexp, code string, tokenType string) Token {
	if m := r.FindString(code); m != "" {
		return Token{tokenType, m, len(m), true}
	}

	return Token{"", "", 0, false}
}

func findSubmatch(r *regexp.Regexp, code string, tokenType string) Token {
	if m := r.FindStringSubmatch(code); len(m) > 0 {
		return Token{tokenType, m[1], len(m[1]), true}
	}

	return Token{"", "", 0, false}
}
