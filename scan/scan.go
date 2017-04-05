//go:generate stringer -type=Type

package scan

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Type identifies the lexical token types of Lisp data.
type Type int

const (
	ILLEGAL Type = iota
	EOF

	ATOM
	LPAREN
	RPAREN
)

// Token represents a token or literal
type Token struct {
	Type Type
	Text string
}

func (t Token) String() string {
	switch t.Type {
	case ATOM:
		return fmt.Sprintf("%s: %q", t.Type, t.Text)
	}
	return t.Type.String()
}

// eofRune is a marker used to indicate that EOF has been reached.
const eofRune = -1

// A Scanner implements reading of Lisp data from an io.Reader.
type Scanner struct {
	Error func(s *Scanner, msg string)

	r  *bufio.Reader
	ch rune
}

func (s *Scanner) error(msg string) {
	if s.Error != nil {
		s.Error(s, msg)
	} else {
		fmt.Fprintf(os.Stderr, "scanner: %s\n", msg)
	}
}

func (s *Scanner) next() {
	r, _, err := s.r.ReadRune()
	switch {
	case err == io.EOF:
		s.ch = eofRune
	case err != nil:
		s.error(err.Error())
	default:
		s.ch = r
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch != eofRune && unicode.IsSpace(s.ch) {
		s.next()
	}
}

func (s *Scanner) scanAtom() string {
	var runes []rune
	for s.ch != eofRune && s.ch != '(' && s.ch != ')' && !unicode.IsSpace(s.ch) {
		runes = append(runes, s.ch)
		s.next()
	}
	return string(runes)
}

// New initialises a scanner for tokenizing Lisp data from a reader.
func New(r io.Reader) *Scanner {
	s := &Scanner{r: bufio.NewReader(r)}
	s.next()
	return s
}

// Next reads the next token from the underlying reader. When the
// token is an ATOM, a literal value is also returned.
//
// In the event of an error, then the Error function will be called if
// it is not nil. Otherwise, an error message is written to Stderr.
func (s *Scanner) Next() Token {
	s.skipWhitespace()
	switch s.ch {
	case eofRune:
		return Token{Type: EOF}
	case '(':
		s.next()
		return Token{Type: LPAREN}
	case ')':
		s.next()
		return Token{Type: RPAREN}
	default:
		return Token{Type: ATOM, Text: s.scanAtom()}
	}
}
