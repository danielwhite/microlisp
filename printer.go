package main

import (
	"bytes"
	"io"
)

// Fprint prints an AST node to w.
func Fprint(w io.Writer, node Node) error {
	var p printer
	p.printNode(node)
	_, err := io.Copy(w, &p.buf)
	return err
}

type printer struct {
	buf bytes.Buffer
}

func (p *printer) printNode(node Node) {
	switch x := node.(type) {
	case *AtomExpr:
		p.printAtom(x)
	case *ListExpr:
		p.printList(x)
	}
}

func (p *printer) printAtom(v *AtomExpr) {
	p.buf.WriteString(v.Name)
}

func (p *printer) printList(v *ListExpr) {
	if v == NIL {
		p.buf.WriteString("()")
		return
	}

	p.buf.WriteByte('(')
	for {
		p.printNode(v.Car)

		next, ok := v.Cdr.(*ListExpr)
		if !ok {
			p.buf.WriteString(" . ")
			p.printNode(v.Cdr)
			break
		}
		if next == NIL {
			break
		}

		p.buf.WriteByte(' ')

		v = next
	}
	p.buf.WriteByte(')')
}
