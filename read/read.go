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
	head := &value.Cell{}
	tail := head
	for {
		tok := r.scanner.Next()
		switch tok.Type {
		case scan.EOF:
			return nil, errEOF
		case scan.Error:
			return nil, errors.New(tok.Text)
		case scan.Comment:
			continue // comments are skipped
		case scan.RightParen:
			tail.Cdr = value.NIL
			return head.Cdr, nil
		case scan.Atom:
			atom := value.Intern(tok.Text)
			cell := &value.Cell{Car: atom}
			tail.Cdr = cell
			tail = cell
		case scan.LeftParen:
			list, err := r.readList()
			if err != nil {
				return nil, err
			}
			cell := &value.Cell{Car: list}
			tail.Cdr = cell
			tail = cell
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
