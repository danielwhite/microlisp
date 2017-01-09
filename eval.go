package main

import "fmt"

type Error string

func (e Error) Error() string {
	return string(e)
}

var DefaultEvaluator = &evaluator{
	env: map[string]Node{
		"t":    T,
		"nil":  NIL,
		"atom": arg1("atom", atom),
		"car":  arg1("car", car),
		"cdr":  arg1("cdr", cdr),
		"cons": arg2("cons", cons),
		"list": Func(list),
	},
}

// Eval returns the value for the given expression within the scope of
// the default evaluator.
func Eval(expr Node) (Node, error) {
	return DefaultEvaluator.Eval(expr)
}

type evaluator struct {
	env map[string]Node
}

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
		if v, ok := e.env[x.Name]; ok {
			return v
		}
		return x // auto-quote
	case *ListExpr:
		switch y := x.Car.(type) {
		case *AtomExpr:
			switch y.Name {
			case "quote":
				return e.evquote(x)
			}
		}
		fn := e.eval(x.Car)
		args := e.evlis(x.Cdr)
		return e.invoke(fn, args)
	default:
		e.errorf("eval: unexpected node of %T: %v", x, x)
	}
	return NIL
}

// invoke applies a list of arguments to a function.
func (e *evaluator) invoke(node Node, args []Node) Node {
	fn, ok := node.(Func)
	if !ok {
		e.errorf("invoke: %s is not a function", node)
	}
	return fn(args)
}

func (e *evaluator) evlis(node Node) []Node {
	if node == NIL {
		return []Node{}
	}

	list, ok := node.(*ListExpr)
	if !ok {
		e.error("evlis: improper argument list")
	}

	var nodes []Node
	next := list
	for {
		nodes = append(nodes, e.eval(next.Car))

		if next.Cdr == NIL {
			break
		}

		v, ok := next.Cdr.(*ListExpr)
		if !ok {
			e.error("evlis: improper argument list")
		}

		next = v
	}
	return nodes
}

// evquote evaluates the quote special form.
func (e *evaluator) evquote(expr *ListExpr) Node {
	v, ok := expr.Cdr.(*ListExpr)
	if !ok || v.Cdr != NIL {
		e.errorf("ill-formed special form: %s", expr)
	}
	return v.Car
}
