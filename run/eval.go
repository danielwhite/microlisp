// Package run provides a Lisp runtime.
package run

import (
	"strings"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

// Eval applies rules to an expression, and returns an expression that
// is the value.
func Eval(expr value.Value) (v value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(value.Error)
		}
	}()
	v = eval(expr, UserEnvironment)
	return
}

// EvalString evaluates the first Lisp expression in a string.
func EvalString(expr string) value.Value {
	scanner := scan.New(strings.NewReader(expr))
	reader := read.New(scanner)

	// Read the next expression from the input.
	v := reader.Read()
	if v == value.EOF {
		return nil
	}
	if err, ok := v.(value.Error); ok {
		return err
	}

	// Evaluate the expression.
	return Eval(v)
}

var (
	// shortcuts for standard types
	T   = value.T
	NIL = value.NIL

	unspecified = value.Error("#[unspecified return value]")
	unassigned  = value.Error("#[unassigned]")
)

func eval(expr value.Value, env value.Environment) value.Value {
	switch x := expr.(type) {
	case *value.Atom:
		if x == T {
			return T
		}
		if x == NIL {
			return NIL
		}
		if v, ok := env.Lookup(x.Name); ok {
			return v
		}
		return x // self-quote unassigned variables
	case *value.Cell:
		return evalForm(x, env)
	}

	panic(expr)
}

func evalForm(expr *value.Cell, env value.Environment) value.Value {
	switch car := expr.Car.(type) {
	case *value.Atom:
		switch car.Name {
		case "quote":
			return evalQuote(expr)
		case "cond":
			return evalCond(expr, env)
		case "lambda":
			return evalLambda(expr, env)
		case "label":
			return evalLabel(expr, env)
		case "defun":
			return evalDefun(expr, env)
		}
	}

	fn := eval(expr.Car, env)

	// Evaluate each argument for application to the function.
	var args []value.Value
	if expr.Cdr != NIL {
		cdr, ok := expr.Cdr.(*value.Cell)
		if !ok {
			value.Errorf("The object %s is not a list", expr.Cdr)
		}
		cdr.Walk(func(v value.Value) {
			args = append(args, eval(v, env))
		})
	}

	return invoke(fn, args)
}

// invoke applies a list of arguments to a function.
func invoke(v value.Value, args []value.Value) value.Value {
	fn, ok := v.(value.Function)
	if !ok {
		value.Errorf("invoke: %s is not a function", v)
	}
	return fn.Invoke(args)
}

// evalQuote evaluates the quote special form.
func evalQuote(expr *value.Cell) value.Value {
	cdr, ok := expr.Cdr.(*value.Cell)
	if !ok || cdr.Cdr != NIL {
		value.Errorf("ill-formed special form: %s", expr)
	}

	return cdr.Car
}

// evalCond evaluates the cond special form.
func evalCond(expr *value.Cell, env value.Environment) value.Value {
	checkExpr := func(ok bool) {
		if !ok {
			value.Errorf("ill-formed special form: %s", expr)
		}
	}

	next := expr
	for next.Cdr != NIL {
		v, ok := next.Cdr.(*value.Cell)
		checkExpr(ok)
		next = v

		clause, ok := next.Car.(*value.Cell)
		checkExpr(ok)

		// if caadr is true, then we want to return the
		// evaluation of the cdadr
		if eval(clause.Car, env) == T {
			body, ok := clause.Cdr.(*value.Cell)
			checkExpr(ok)
			return evalProgn(body, env)
		}
	}
	return NIL
}

// evalProgn evaluates a proper list of values, returning the last
// value.
func evalProgn(v value.Value, env value.Environment) value.Value {
	var last value.Value = unspecified
	if v == NIL {
		return last
	}

	list, ok := v.(*value.Cell)
	if !ok {
		value.Errorf("implicit progn must be a list: %s", v)
	}

	list.Walk(func(v value.Value) {
		last = eval(v, env)
	})
	return last
}

func evalLambda(expr *value.Cell, env value.Environment) value.Value {
	checkExpr := func(ok bool) {
		if !ok {
			value.Errorf("ill-formed special form: %s", expr)
		}
	}

	// (cadr (lambda (arg1 ... argN) body1 ... bodyN))
	cdr, ok := expr.Cdr.(*value.Cell)
	checkExpr(ok && cdr.Cdr != NIL)

	// (cddr (lambda args body1 body2 ... bodyN))
	body, ok := cdr.Cdr.(*value.Cell)
	checkExpr(ok)

	return makeFunction(cdr.Car, body, env)
}

// evalLabel evaluates the label special form.
func evalLabel(expr *value.Cell, env value.Environment) value.Value {
	checkExpr := func(ok bool) {
		if !ok {
			value.Errorf("ill-formed special form: %s", expr)
		}
	}

	// (cadr (label name (lambda ...)))
	cdr, ok := expr.Cdr.(*value.Cell)
	checkExpr(ok && cdr.Cdr != NIL)
	label, ok := cdr.Car.(*value.Atom)
	checkExpr(ok)

	// (caddr (label name (lambda ...)))
	cddr, ok := cdr.Cdr.(*value.Cell)
	checkExpr(ok && cddr.Cdr == NIL)
	caddr, ok := cddr.Car.(*value.Cell)
	checkExpr(ok)

	// Evaluate lambda in an environment where it is able to
	// reference the name defined by the label special form.
	extEnv := value.NewEnv(env)
	extEnv.Define(label.Name, unassigned)
	fn := evalLambda(caddr, extEnv)

	// Update the binding to the newly created function.
	extEnv.Update(label.Name, fn)

	return fn
}

// evalDefun defines a function permanently.
func evalDefun(expr *value.Cell, env value.Environment) *value.Atom {
	checkExpr := func(ok bool) {
		if !ok {
			value.Errorf("ill-formed special form: %s", expr)
		}
	}

	// (cadr (defun fn (arg1 ... argN) body1 ... bodyN))
	cdr, ok := expr.Cdr.(*value.Cell)
	checkExpr(ok && cdr.Cdr != NIL)
	symbol, ok := cdr.Car.(*value.Atom)
	checkExpr(ok)

	// (caddr (defun fn (arg1 ... argN) body1 ... bodyN))
	cddr, ok := cdr.Cdr.(*value.Cell)
	checkExpr(ok && cddr.Cdr != NIL)

	// (cdddr (defun fn (arg1 ... argN) body1 ... bodyN))
	body, ok := cddr.Cdr.(*value.Cell)
	checkExpr(ok)

	// By defining in the current environment, we add a permanent
	// function, but don't need to find the toplevel. I think this
	// differs from a typical Lisp, but falls within McCarthy's
	// described behaviour.
	env.Define(symbol.Name, unassigned)
	fn := makeFunction(cddr.Car, body, env)
	env.Update(symbol.Name, fn)

	return symbol
}

// makeFunction creates a new function from the lambda special form.
func makeFunction(argExpr value.Value, bodyExpr *value.Cell, env value.Environment) value.Function {
	var vars []string
	if argExpr != NIL {
		argCell, ok := argExpr.(*value.Cell)
		if !ok {
			value.Errorf("The object %s is not a list", argExpr)
		}

		argCell.Walk(func(v value.Value) {
			atom, ok := v.(*value.Atom)
			if !ok {
				value.Errorf("The object %s is not a symbol", v)
			}

			vars = append(vars, atom.Name)
		})
	}

	fn := func(args []value.Value) value.Value {
		// Arguments are bound in a new environment.
		extEnv := value.NewEnv(env)
		for i, name := range vars {
			extEnv.Define(name, args[i])
		}

		// Body is evaluated in an implicit progn.
		var last value.Value = unspecified
		bodyExpr.Walk(func(v value.Value) {
			last = eval(v, extEnv)
		})
		return last
	}

	return value.FuncX(len(vars), fn)
}
