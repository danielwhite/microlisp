//go:generate stringer -type=Type

package scan

import (
	"fmt"
	"io"
	"unicode"
)

// Type identifies the lexical token types of Lisp data.
type Type int

const (
	Illegal Type = iota
	Error
	EOF

	Atom
	LeftParen
	RightParen
)

// Token represents a token or literal
type Token struct {
	Type Type
	Text string
}

func (t Token) String() string {
	switch t.Type {
	case Error:
		return fmt.Sprintf("error: %s", t.Text)
	case Atom:
		return fmt.Sprintf("%s: %q", t.Type, t.Text)
	}
	return t.Type.String()
}

// Scanner holds state of Lisp tokens.
type Scanner struct {
	r   io.RuneReader // The reader provided by the client.
	ch  rune          // Last character read.
	err error         // Sticky error.
}

// marker used to indicate that EOF has been reached
const eof = -1

func (s *Scanner) readChar() {
	s.ch, _, s.err = s.r.ReadRune()
	if s.err != nil {
		s.ch = eof
	}
}

func (s *Scanner) lexAtom() Token {
	var text []rune
	for s.ch != eof && s.ch != '(' && s.ch != ')' && !unicode.IsSpace(s.ch) {
		text = append(text, s.ch)
		s.readChar()
	}
	return Token{Type: Atom, Text: string(text)}
}

// New initialises a scanner for tokenizing Lisp data from a reader.
func New(r io.RuneReader) *Scanner {
	return &Scanner{
		r:  r,
		ch: ' ', // start with whitespace that is dropped
	}
}

// Next reads the next token from the underlying reader.
//
// If an error is encountered, an error token will be returned with a
// message as its text.
func (s *Scanner) Next() Token {
	// All whitespace is ignored.
	for unicode.IsSpace(s.ch) {
		s.readChar()
	}

	switch s.ch {
	case '(':
		s.readChar()
		return Token{Type: LeftParen}
	case ')':
		s.readChar()
		return Token{Type: RightParen}
	case eof:
		if s.err == io.EOF {
			return Token{Type: EOF}
		}
		return Token{Type: Error, Text: s.err.Error()}
	default:
		return s.lexAtom()
	}
}
