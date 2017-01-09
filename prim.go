package main

import "fmt"

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

func arg1(fn func(Node) Node) Func {
	return func(args []Node) Node {
		if len(args) != 1 {
			errorf("car: called with %d arguments; requires exactly 1 argument", len(args))
		}
		return fn(args[0])
	}
}

func errorf(format string, a ...interface{}) {
	panic(Error(fmt.Sprintf(format, a...)))
}
