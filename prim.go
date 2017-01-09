package main

import (
	"fmt"
	"reflect"
)

func atom(arg Node) Node {
	if _, ok := arg.(*AtomExpr); ok {
		return T
	}
	return NIL
}

func equal(a Node, b Node) Node {
	if reflect.DeepEqual(a, b) {
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
	return func(args []Node) Node {
		if len(args) != 1 {
			errorf("%s: called with %d arguments; requires exactly 1 argument", name, len(args))
		}
		return fn(args[0])
	}
}

func arg2(name string, fn func(Node, Node) Node) Func {
	return func(args []Node) Node {
		if len(args) != 2 {
			errorf("%s: called with %d arguments; requires exactly 2 argument", name, len(args))
		}
		return fn(args[0], args[1])
	}
}

func errorf(format string, a ...interface{}) {
	panic(Error(fmt.Sprintf(format, a...)))
}
