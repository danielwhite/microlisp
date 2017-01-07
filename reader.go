package main

import (
	"errors"
	"fmt"
	"io"
)

// Read parses the next expression from a stream of tokens. When the
// end of the stream is reached, then io.EOF is returned.
func Read(s *Scanner) (node Node, err error) {
	tok, lit := s.Scan()
	switch tok {
	case ILLEGAL:
		err = errors.New("illegal token")
	case ATOM:
		node = &AtomExpr{Name: lit}
	case LPAREN:
		node, err = readList(s)
	case RPAREN:
		err = errors.New("unbalanced closed parenthesis")
	case EOF:
		err = io.EOF
	default:
		err = fmt.Errorf("unsupported token: %s: %q", tok, lit)
	}
	return
}

func readList(s *Scanner) (v *ListExpr, err error) {
	tok, lit := s.Scan()
	switch tok {
	case EOF:
		err = errors.New("premature EOF")
	case ATOM:
		v = &ListExpr{Car: &AtomExpr{Name: lit}}
		v.Cdr, err = readList(s)
	case LPAREN:
		v = &ListExpr{}
		v.Car, err = readList(s)
		if err == nil {
			v.Cdr, err = readList(s)
		}
	case RPAREN:
		// nothing to do
	default:
		err = fmt.Errorf("unexpected token in list: %s", tok)
	}
	return
}
