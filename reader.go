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
	case scan.Illegal:
		err = errors.New("illegal token")
	case scan.Error:
		err = errors.New(tok.Text)
	case scan.Atom:
		node = &AtomExpr{Name: tok.Text}
	case scan.LeftParen:
		node, err = readList(s)
	case scan.RightParen:
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
	case scan.Atom:
		var cdr Node
		cdr, err = readList(s)
		if err != nil {
			break
		}
		node = &ListExpr{&AtomExpr{Name: tok.Text}, cdr}
	case scan.LeftParen:
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
	case scan.RightParen:
		node = NIL
	default:
		err = fmt.Errorf("unexpected token in list: %s", tok)
	}
	return
}
