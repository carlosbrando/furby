// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package new_lexer

import "fmt"

// itemType identifies the type of lex items.
type itemType int

const (
	itemError      itemType = iota // error ocurred; value is text of error
	itemDot                        // the cursor, spelled '.'
	itemEOF                        //
	itemElse                       // else keyword
	itemEnd                        // end keyword
	itemField                      // identifier, starting with '.'
	itemIdentifier                 // identifier
	itemIf                         // if keyword
	itemLeftMeta                   // left meta-string
	itemNumber                     // number
	itemPipe                       // pipe symbol
	itemRange                      // range keyword
	itemRawString                  // raw quoted string (includes quotes)
	itemRightMeta                  // right meta-string
	itemString                     // quoted string (includes quotes)
	itemText                       // plain text
)

// item represents a token returned from the scanner.
type item struct {
	typ itemType // Type, such a itemNumber.
	val string   // Value, such as "23.2".
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}

	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}
