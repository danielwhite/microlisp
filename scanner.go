package main

//go:generate stringer -type=Token

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Token is the set of lexical tokens of Lisp data.
type Token int

const (
	ILLEGAL Token = iota
	EOF

	ATOM
	LPAREN
	RPAREN
)

// A Scanner implements reading of Lisp data from an io.Reader.
type Scanner struct {
	Error func(s *Scanner, msg string)

	r  *bufio.Reader
	ch rune
}

// eofRune is a marker used to indicate that EOF has been reached.
const eofRune = -1

// NewScanner returns a scanner for tokenizing Lisp data.
func NewScanner(r io.Reader) *Scanner {
	s := &Scanner{r: bufio.NewReader(r)}
	s.next()
	return s
}

// Scan reads the next token from the underlying reader. When the
// token is an ATOM, a literal value is also returned.
//
// In the event of an error, then the Error function will be called if
// it is not nil. Otherwise, an error message is written to Stderr.
func (s *Scanner) Scan() (tok Token, lit string) {
	s.skipWhitespace()
	switch s.ch {
	case eofRune:
		tok = EOF
	case '(':
		tok = LPAREN
		s.next()
	case ')':
		tok = RPAREN
		s.next()
	default:
		tok = ATOM
		lit = s.scanAtom()
	}
	return
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
