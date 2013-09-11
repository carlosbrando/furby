// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"unicode/utf8"
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
	// itemText 			// plain text
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

const eof = -1

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

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	// go l.run()
	l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	// for l.state = lexText; l.state != nil; {
	for l.state = lexAction; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

// lexAction scans the elements inside action delimiters.
func lexAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	// if strings.HasPrefix(l.input[l.pos:], l.rightDelim) {
	// 	if l.parenDepth == 0 {
	// 		return lexRightDelim
	// 	}
	// 	return l.errorf("unclosed left paren")
	// }
	switch r := l.next(); {
	case r == eof: // || isEndOfLine(r):
		return nil
		// return l.errorf("unclosed action")
	// case isSpace(r):
	// 	return lexSpace
	// case r == ':':
	// 	if l.next() != '=' {
	// 		return l.errorf("expected :=")
	// 	}
	// 	l.emit(itemColonEquals)
	// case r == '|':
	// 	l.emit(itemPipe)
	// case r == '"':
	// 	return lexQuote
	// case r == '`':
	// 	return lexRawQuote
	// case r == '$':
	// 	return lexVariable
	// case r == '\'':
	// 	return lexChar
	// case r == '.':
	// 	// special look-ahead for ".field" so we don't break l.backup().
	// 	if l.pos < Pos(len(l.input)) {
	// 		r := l.input[l.pos]
	// 		if r < '0' || '9' < r {
	// 			return lexField
	// 		}
	// 	}
	// 	fallthrough // '.' can start a number.
	// case r == '+' || r == '-' || ('0' <= r && r <= '9'):
	// 	l.backup()
	// 	return lexNumber
	// case isAlphaNumeric(r):
	// 	l.backup()
	// 	return lexIdentifier
	// case r == '(':
	// 	l.emit(itemLeftParen)
	// 	l.parenDepth++
	// 	return lexAction
	// case r == ')':
	// 	l.emit(itemRightParen)
	// 	l.parenDepth--
	// 	if l.parenDepth < 0 {
	// 		return l.errorf("unexpected right paren %#U", r)
	// 	}
	// 	return lexAction
	// case r <= unicode.MaxASCII && unicode.IsPrint(r):
	// 	l.emit(itemChar)
	// 	return lexAction
	default:
		fmt.Println(string(r))
		// return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexAction
}
