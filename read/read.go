// Package read implements a reader for Lisp expressions.
package read

import (
	"errors"
	"fmt"

	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

var errEOF = errors.New("premature EOF")

// Reader holds state of Lisp data.
type Reader struct {
	scanner *scan.Scanner
}

func (r *Reader) readList() (value.Value, error) {
	var list value.List
	for {
		switch tok := r.scanner.Next(); tok.Type {
		case scan.Atom:
			list = append(list, value.Intern(tok.Text))
		case scan.LeftParen:
			v, err := r.readList()
			if err != nil {
				return nil, err
			}
			list = append(list, v)
		case scan.RightParen:
			if len(list) == 0 {
				return value.NIL, nil
			}
			return append(list, value.NIL), nil
		case scan.EOF:
			return nil, errEOF
		case scan.Error:
			return nil, errors.New(tok.Text)
		case scan.Comment:
			// comments are ignored
		default:
			return nil, fmt.Errorf("unexpected token in list: %s", tok)
		}
	}
}

// New initialises a reader for parsing Lisp expressions.
func New(s *scan.Scanner) *Reader {
	return &Reader{scanner: s}
}

// Read parses the next expression from a stream of tokens. When the
// end of the stream is reached, then value.EOF.
func (r *Reader) Read() value.Value {
Next:
	switch tok := r.scanner.Next(); tok.Type {
	case scan.Atom:
		return value.Intern(tok.Text)
	case scan.LeftParen:
		v, err := r.readList()
		if err != nil {
			return value.Error(err.Error())
		}
		return v
	case scan.RightParen:
		return value.Error("unbalanced closed parenthesis")
	case scan.EOF:
		return value.EOF
	case scan.Error:
		return value.Error(tok.Text)
	case scan.Comment:
		// comments are ignored, so scan again
		goto Next
	default:
		return value.Error(fmt.Sprintf("unsupported token: %s", tok))
	}
}
