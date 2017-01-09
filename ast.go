package main

type Node interface{}

// A List node represents a cell in a list.
type ListExpr struct {
	Car Node
	Cdr Node
}

func (e *ListExpr) String() string {
	return Sprint(e)
}

// NIL is the empty list.
var NIL = (*ListExpr)(nil)

type AtomExpr struct {
	Name string
}

func (e *AtomExpr) String() string {
	return e.Name
}

// Func represents a function that can be invoked.
type Func func([]Node) Node
