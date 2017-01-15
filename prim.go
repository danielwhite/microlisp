package main

import "fmt"

func atom(arg Node) Node {
	if _, ok := arg.(*AtomExpr); ok {
		return T
	}
	return NIL
}

func equal(a Node, b Node) Node {
	if a.Equal(b) {
		return T
	}
	return NIL
}

func car(arg Node) Node {
	v, ok := arg.(*ListExpr)
	if !ok {
		errorf("car: %v is not a pair", arg)
	}
	return v.Car
}

func cdr(arg Node) Node {
	v, ok := arg.(*ListExpr)
	if !ok {
		errorf("cdr: %v is not a pair", arg)
	}
	return v.Cdr
}

func cons(a Node, b Node) Node {
	return &ListExpr{a, b}
}

func list(args []Node) Node {
	if len(args) == 0 {
		return NIL
	}

	list := &ListExpr{args[0], NIL}
	next := list
	for _, arg := range args[1:] {
		cons := &ListExpr{arg, NIL}
		next.Cdr = cons
		next = cons
	}
	return list
}

func arg1(name string, fn func(Node) Node) Func {
	return &Func1{fn}
}

func arg2(name string, fn func(Node, Node) Node) Func {
	return &Func2{fn}
}

func argN(name string, fn func([]Node) Node) Func {
	return &FuncN{fn}
}

func errorf(format string, a ...interface{}) {
	panic(Error(fmt.Sprintf(format, a...)))
}
