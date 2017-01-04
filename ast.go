package main

type Node interface{}

// A List node represents a
type ListExpr struct {
	Car Node
	Cdr Node
}

// NIL is the empty list.
var NIL = (*ListExpr)(nil)

type AtomExpr struct {
	Name string
}
