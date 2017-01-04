package main

import (
	"bytes"
	"errors"
	"fmt"
)

// Parse Lisp data into an abstract syntax tree.
func Parse(src []byte) (node Node, err error) {
	s := NewScanner(bytes.NewReader(src))
	tok, lit := s.Scan()
	switch tok {
	case ILLEGAL:
		err = errors.New("illegal token")
	case ATOM:
		node = &AtomExpr{Name: lit}
	case LPAREN:
		node, err = parseList(s)
	case RPAREN:
		err = errors.New("Unbalanced closed parenthesis")
	case EOF:
		// nothing to do
	default:
		err = fmt.Errorf("unsupported token: %s: %q", tok, lit)
	}
	return
}

func parseList(s *Scanner) (v *ListExpr, err error) {
	tok, lit := s.Scan()
	switch tok {
	case EOF:
		err = errors.New("Premature EOF")
	case ATOM:
		v = &ListExpr{Car: &AtomExpr{Name: lit}}
		v.Cdr, err = parseList(s)
	case LPAREN:
		v = &ListExpr{}
		v.Car, err = parseList(s)
		if err == nil {
			v.Cdr, err = parseList(s)
		}
	case RPAREN:
		// nothing to do
	default:
		err = fmt.Errorf("unexpected token in list: %s", tok)
	}
	return
}
