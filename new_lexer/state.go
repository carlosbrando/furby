// Copyright 2013 Carlos Brando. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package new_lexer

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn
