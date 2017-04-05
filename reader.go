package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/danielwhite/microlisp/scan"
)

// Read parses the next expression from a stream of tokens. When the
// end of the stream is reached, then io.EOF is returned.
func Read(s *scan.Scanner) (node Node, err error) {
	tok := s.Next()
	switch tok.Type {
	case scan.ILLEGAL:
		err = errors.New("illegal token")
	case scan.ATOM:
		node = &AtomExpr{Name: tok.Text}
	case scan.LPAREN:
		node, err = readList(s)
	case scan.RPAREN:
		err = errors.New("unbalanced closed parenthesis")
	case scan.EOF:
		err = io.EOF
	default:
		err = fmt.Errorf("unsupported token: %s", tok)
	}
	return
}

func readList(s *scan.Scanner) (node Node, err error) {
	tok := s.Next()
	switch tok.Type {
	case scan.EOF:
		err = errors.New("premature EOF")
	case scan.ATOM:
		var cdr Node
		cdr, err = readList(s)
		if err != nil {
			break
		}
		node = &ListExpr{&AtomExpr{Name: tok.Text}, cdr}
	case scan.LPAREN:
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
	case scan.RPAREN:
		node = NIL
	default:
		err = fmt.Errorf("unexpected token in list: %s", tok)
	}
	return
}
