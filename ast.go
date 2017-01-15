package main

type Node interface {
	Equal(Node) bool
}

// A List node represents a cell in a list.
type ListExpr struct {
	Car Node
	Cdr Node
}

func (e *ListExpr) Equal(node Node) bool {
	v, ok := node.(*ListExpr)
	if ok && e.Car.Equal(v.Car) && e.Cdr.Equal(v.Cdr) {
		return true
	}
	return false
}

func (e *ListExpr) String() string {
	return Sprint(e)
}

// NIL is the empty list.
var NIL = &AtomExpr{"nil"}

// T is the true value.
var T = &AtomExpr{"t"}

type AtomExpr struct {
	Name string
}

func (e *AtomExpr) Equal(node Node) bool {
	v, ok := node.(*AtomExpr)
	if !ok {
		return false
	}

	return e.Name == v.Name
}

func (e *AtomExpr) String() string {
	return e.Name
}

// Func represents a function that can be invoked.
type Func interface {
	Node
	Invoke([]Node) Node
}

// Func1 represents a function that accepts a single argument.
type Func1 struct {
	fn func(Node) Node
}

func (f *Func1) Invoke(args []Node) Node {
	if len(args) != 1 {
		errorf("called with %d arguments; requires exactly 1 argument", len(args))
	}
	return f.fn(args[0])
}

func (f *Func1) Equal(node Node) bool {
	return f == node
}

// Func2 represents a function that accepts exactly two arguments.
type Func2 struct {
	fn func(Node, Node) Node
}

func (f *Func2) Invoke(args []Node) Node {
	if len(args) != 2 {
		errorf("called with %d arguments; requires exactly 1 argument", len(args))
	}
	return f.fn(args[0], args[1])
}

func (f *Func2) Equal(node Node) bool {
	return f == node
}

// FuncN represents a function that accepts a variable number
// arugments.
type FuncN struct {
	fn func([]Node) Node
}

func (f *FuncN) Invoke(args []Node) Node {
	return f.fn(args)
}

func (f *FuncN) Equal(node Node) bool {
	return f == node
}
