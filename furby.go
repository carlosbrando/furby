// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/carlosbrando/furby/lexer"
	"github.com/carlosbrando/furby/parse"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// read code
	code, err := ioutil.ReadFile("hello.frb")
	check(err)

	tokens := lexer.Tokenize(string(code))
	fmt.Println(tokens)

	// The all new code.
	fmt.Println("\n----------")
	parse.Parse("hello.frb", string(code))
}
