// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

import (
	"fmt"
)

type Tree struct {
}

func Parse(name, text string) {
	x := lex(name, text)
	fmt.Println(x.input)
}
