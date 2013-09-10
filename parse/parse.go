// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

type Tree struct {
}

func Parse(name, text string) {
	lex(name, text)
}
