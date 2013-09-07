package lexer

import "regexp"

func find(r *regexp.Regexp, code string, tokenType string) Token {
	if m := r.FindString(code); m != "" {
		return Token{tokenType, m, len(m), true}
	}

	return Token{"", "", 0, false}
}
