// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"strings"
	"unicode"
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
	itemBool                  // boolean constant
	// itemChar                         // printable ASCII character; grab bag for comma etc.
	// itemCharConstant                 // character constant
	itemComplex // complex constant (1+2i); imaginary is just a number
	// itemColonEquals                  // colon-equals (':=') introducing a declaration
	itemEOF
	itemField      // alphanumeric identifier starting with '.'
	itemIdentifier // alphanumeric identifier not starting with '.'
	// itemLeftDelim  // left action delimiter
	// itemLeftParen  // '(' inside action
	itemNumber // simple number, including imaginary
	// itemPipe       // pipe symbol
	// itemRawString  // raw quoted string (includes quotes)
	// itemRightDelim // right action delimiter
	// itemRightParen // ')' inside action
	itemSpace // run of spaces separating arguments
	itemEndOfLine
	// itemString     // quoted string (includes quotes)
	// itemText 			// plain text
	// itemVariable   // variable starting with '$', such as '$' or  '$1' or '$hello'
	// Keywords appear after all the rest.
	itemKeyword // used only to delimit the keywords
	// itemDot      // the cursor, spelled '.'
	// itemDefine   // define keyword
	itemElse // else keyword
	itemEnd  // end keyword
	itemIf   // if keyword
	itemNil  // the untyped nil constant, easiest to treat as a keyword
	// itemRange    // range keyword
	// itemTemplate // template keyword
	// itemWith     // with keyword
)

var key = map[string]itemType{
	// ".": itemDot,
	// "define":   itemDefine,
	"else": itemElse,
	"end":  itemEnd,
	"if":   itemIf,
	// "range":    itemRange,
	"nil": itemNil,
	// "template": itemTemplate,
	// "with": itemWith,
}

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

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	// fmt.Println(item{t, l.start, l.input[l.start:l.pos]})
	// fmt.Println("----------")
	// TODO: Remove above and leave below
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	// fmt.Println(item{itemError, l.start, fmt.Sprintf(format, args...)})
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	// l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexAction; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

// lexAction scans all the elements.
func lexAction(l *lexer) stateFn {
	// Either number, quoted string, or identifier.
	// Spaces separate arguments; runs of spaces turn into itemSpace.
	// Pipe symbols separate and are emitted.
	switch r := l.next(); {
	case r == eof:
		return nil
	case isEndOfLine(r):
		l.emit(itemEndOfLine)
	case isSpace(r):
		return lexSpace
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
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
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
		// fmt.Println(string(r))
		// return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexAction
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexAction
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(itemComplex)
	} else {
		l.emit(itemNumber)
	}
	return lexAction
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > itemKeyword:
				l.emit(key[word])
			case word[0] == '.':
				l.emit(itemField)
			case word == "true", word == "false":
				l.emit(itemBool)
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}
	return lexAction
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(':
		return true
	}
	return false
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
