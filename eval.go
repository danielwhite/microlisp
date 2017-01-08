package main

import "fmt"

type Error string

func (e Error) Error() string {
	return string(e)
}

var DefaultEvaluator = &evaluator{}

// Eval returns the value for the given expression within the scope of
// the default evaluator.
func Eval(expr Node) (Node, error) {
	return DefaultEvaluator.Eval(expr)
}

type evaluator struct{}

func (e *evaluator) error(err string) {
	panic(Error(err))
}

func (e *evaluator) errorf(format string, a ...interface{}) {
	panic(Error(fmt.Sprintf(format, a...)))
}

// Eval applies rules to an expression, and returns an expression that
// is the value.
func (e *evaluator) Eval(expr Node) (value Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(Error)
		}
	}()
	value = e.eval(expr)
	return
}

func (e *evaluator) eval(node Node) Node {
	switch x := node.(type) {
	case *AtomExpr:
		return x // auto-quote
	case *ListExpr:
		switch y := x.Car.(type) {
		case *AtomExpr:
			switch y.Name {
			case "quote":
				return e.cadr(x)
			default:
				e.errorf("eval: %q is an unknown special form", y.Name)
			}
		}
	}
	return NIL
}

func (e *evaluator) cadr(v *ListExpr) Node {
	v, ok := v.Cdr.(*ListExpr)
	if !ok || v == NIL {
		e.errorf("cadr: %v is not a pair", v)
	}
	return v.Car
}
