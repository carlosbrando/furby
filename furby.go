package main

import (
	"fmt"
	"io/ioutil"

	"github.com/carlosbrando/furby/lexer"
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
}
