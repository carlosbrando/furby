// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

import (
	"fmt"
)

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos Pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	// itemBool                         // boolean constant
	// itemChar                         // printable ASCII character; grab bag for comma etc.
	// itemCharConstant                 // character constant
	// itemComplex                      // complex constant (1+2i); imaginary is just a number
	// itemColonEquals                  // colon-equals (':=') introducing a declaration
	itemEOF
	// itemField      // alphanumeric identifier starting with '.'
	// itemIdentifier // alphanumeric identifier not starting with '.'
	// itemLeftDelim  // left action delimiter
	// itemLeftParen  // '(' inside action
	// itemNumber     // simple number, including imaginary
	// itemPipe       // pipe symbol
	// itemRawString  // raw quoted string (includes quotes)
	// itemRightDelim // right action delimiter
	// itemRightParen // ')' inside action
	// itemSpace      // run of spaces separating arguments
	// itemString     // quoted string (includes quotes)
	// itemText       // plain text
	// itemVariable   // variable starting with '$', such as '$' or  '$1' or '$hello'
	// Keywords appear after all the rest.
	itemKeyword // used only to delimit the keywords
	// itemDot      // the cursor, spelled '.'
	// itemDefine   // define keyword
	// itemElse     // else keyword
	// itemEnd      // end keyword
	// itemIf       // if keyword
	// itemNil      // the untyped nil constant, easiest to treat as a keyword
	// itemRange    // range keyword
	// itemTemplate // template keyword
	// itemWith     // with keyword
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	lastPos    Pos       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	// go l.run()
	return l
}

// run runs the state machine for the lexer.
// func (l *lexer) run() {
// 	for l.state = lexText; l.state != nil; {
// 		l.state = l.state(l)
// 	}
// }
