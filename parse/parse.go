// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"runtime"
)

// Tree is the representation of a single parsed template.
type Tree struct {
	Name      string    // name of the template represented by the tree.
	ParseName string    // name of the top-level template during parsing, for error messages.
	Root      *ListNode // top-level root of the tree.
	text      string    // text parsed to create the template (or its parent)
	// Parsing only; cleared after parse.
	funcs     []map[string]interface{}
	lex       *lexer
	token     [3]item // three-token lookahead for parser.
	peekCount int
	vars      []string // variables defined at the moment.
}

// Parse returns a map from template name to parse.Tree, created by parsing the
// templates described in the argument string. The top-level template will be
// given the specified name. If an error is encountered, parsing stops and an
// empty map is returned with the error.
func Parse(name, text string, funcs ...map[string]interface{}) (treeSet map[string]*Tree, err error) {
	treeSet = make(map[string]*Tree)
	t := New(name)
	t.text = text
	_, err = t.Parse(text, treeSet, funcs...)
	return
}

// Parse parses the template definition string to construct a representation of
// the template for execution. If either action delimiter string is empty, the
// default ("{{" or "}}") is used. Embedded template definitions are added to
// the treeSet map.
func (t *Tree) Parse(text string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error) {
	defer t.recover(&err)
	t.ParseName = t.Name
	t.startParse(funcs, lex(t.Name, text))
	t.text = text
	t.parse(treeSet)
	t.add(treeSet)
	t.stopParse()
	return t, nil
}

// add adds tree to the treeSet.
func (t *Tree) add(treeSet map[string]*Tree) {
	tree := treeSet[t.Name]
	if tree == nil || IsEmptyTree(tree.Root) {
		treeSet[t.Name] = t
		return
	}
	if !IsEmptyTree(t.Root) {
		t.errorf("template: multiple definition of template %q", t.Name)
	}
}

// IsEmptyTree reports whether this tree (node) is empty of everything but space.
func IsEmptyTree(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	case *ActionNode:
	// case *IfNode:
	case *ListNode:
		for _, node := range n.Nodes {
			if !IsEmptyTree(node) {
				return false
			}
		}
		return true
	// case *RangeNode:
	// case *TemplateNode:
	// case *TextNode:
	// 	return len(bytes.TrimSpace(n.Text)) == 0
	// case *WithNode:
	default:
		panic("unknown node: " + n.String())
	}
	return false
}

// next returns the next token.
func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// peek returns but does not consume the next token.
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

// nextNonSpace returns the next non-space token.
func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

// Parsing.

// New allocates a new parse tree with the given name.
func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("template: %s:%d: %s", t.ParseName, t.lex.lineNumber(), format)
	panic(fmt.Errorf(format, args...))
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}

// startParse initializes the parser, using the lexer.
func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer) {
	t.Root = nil
	t.lex = lex
	t.vars = []string{"$"}
	t.funcs = funcs
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
	t.vars = nil
	t.funcs = nil
}

// parse is the top-level parser for a template, essentially the same
// as itemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) parse(treeSet map[string]*Tree) (next Node) {
	t.Root = newList(t.peek().pos)
	for t.peek().typ != itemEOF {
		// if t.peek().typ == itemLeftDelim {
		// 	delim := t.next()
		// 	if t.nextNonSpace().typ == itemDefine {
		// 		newT := New("definition") // name will be updated once we know it.
		// 		newT.text = t.text
		// 		newT.ParseName = t.ParseName
		// 		newT.startParse(t.funcs, t.lex)
		// 		newT.parseDefinition(treeSet)
		// 		continue
		// 	}
		// 	t.backup2(delim)
		// }
		// n := t.textOrAction()
		n := t.action()
		if n.Type() == nodeEnd {
			t.errorf("unexpected %s", n)
		}
		t.Root.append(n)
	}
	return nil
}

// Action:
//  control
//  command ("|" command)*
// Left delim is past. Now get actions.
// First word could be a keyword such as range.
func (t *Tree) action() (n Node) {
	switch token := t.nextNonSpace(); token.typ {
	// case itemElse:
	// 	return t.elseControl()
	// case itemEnd:
	// 	return t.endControl()
	// case itemIf:
	// 	return t.ifControl()
	// case itemRange:
	// 	return t.rangeControl()
	// case itemTemplate:
	// 	return t.templateControl()
	// case itemWith:
	// 	return t.withControl()
	}
	t.backup()
	// Do not pop variables; they persist until "end".
	return newAction(t.peek().pos, t.lex.lineNumber())
}
