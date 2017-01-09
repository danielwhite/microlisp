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

func readList(s *Scanner) (node Node, err error) {
	tok, lit := s.Scan()
	switch tok {
	case EOF:
		err = errors.New("premature EOF")
	case ATOM:
		var cdr Node
		cdr, err = readList(s)
		if err != nil {
			break
		}
		node = &ListExpr{&AtomExpr{Name: lit}, cdr}
	case LPAREN:
		var car, cdr Node
		car, err = readList(s)
		if err != nil {
			break
		}
		cdr, err = readList(s)
		if err != nil {
			break
		}
		node = &ListExpr{car, cdr}
	case RPAREN:
		node = NIL
	default:
		err = fmt.Errorf("unexpected token in list: %s", tok)
	}
	return
}
