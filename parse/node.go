// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Parse nodes.

package parse

import (
	"bytes"
	"fmt"
)

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
	Type() NodeType
	String() string
	// Copy does a deep copy of the Node and all its components.
	// To avoid type assertions, some XxxNodes also have specialized
	// CopyXxx methods that return *XxxNode.
	Copy() Node
	Position() Pos // byte position of start of node in full original input string
	// Make sure only functions in this package can create Nodes.
	unexported()
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// unexported keeps Node implementations local to the package.
// All implementations embed Pos, so this takes care of it.
func (Pos) unexported() {
}

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText   NodeType = iota // Plain text.
	NodeAction                 // A non-control action such as a field evaluation.
	// NodeBool                       // A boolean constant.
	// NodeChain                      // A sequence of field accesses.
	// NodeCommand                    // An element of a pipeline.
	// NodeDot                        // The cursor, dot.
	// nodeElse                       // An else action. Not added to tree.
	nodeEnd // An end action. Not added to tree.
	// NodeField                      // A field or method name.
	// NodeIdentifier                 // An identifier; always a function name.
	// NodeIf                         // An if action.
	NodeList // A list of Nodes.
	// NodeNil                        // An untyped nil constant.
	// NodeNumber                     // A numerical constant.
	// NodePipe                       // A pipeline of commands.
	// NodeRange                      // A range action.
	// NodeString                     // A string constant.
	// NodeTemplate                   // A template invocation action.
	// NodeVariable                   // A $ variable.
	// NodeWith                       // A with action.
)

// Nodes.

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Pos
	Nodes []Node // The element nodes in lexical order.
}

func newList(pos Pos) *ListNode {
	return &ListNode{NodeType: NodeList, Pos: pos}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) String() string {
	b := new(bytes.Buffer)
	for _, n := range l.Nodes {
		fmt.Fprint(b, n)
	}
	return b.String()
}

func (l *ListNode) CopyList() *ListNode {
	if l == nil {
		return l
	}
	n := newList(l.Pos)
	for _, elem := range l.Nodes {
		n.append(elem.Copy())
	}
	return n
}

func (l *ListNode) Copy() Node {
	return l.CopyList()
}

// ActionNode holds an action (something bounded by delimiters).
// Control actions have their own nodes; ActionNode represents simple
// ones such as field evaluations and parenthesized pipelines.
type ActionNode struct {
	NodeType
	Pos
	Line int // The line number in the input (deprecated; kept for compatibility)
}

func (a *ActionNode) String() string {
	// return fmt.Sprintf("{{%s}}", a.Pipe)
	return ""

}

func (a *ActionNode) Copy() Node {
	return newAction(a.Pos, a.Line)

}

func newAction(pos Pos, line int) *ActionNode {
	return &ActionNode{NodeType: NodeAction, Pos: pos, Line: line}
}

// CommandNode holds a command (a pipeline inside an evaluating action).
type CommandNode struct {
	NodeType
	Pos
	Args []Node // Arguments in lexical order: Identifier, field, or constant.
}

// VariableNode holds a list of variable names, possibly with chained field
// accesses. The dollar sign is part of the (first) name.
type VariableNode struct {
	NodeType
	Pos
	Ident []string // Variable name and fields in lexical order.
}
