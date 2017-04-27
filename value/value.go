// Package value implements Lisp values and their evaluation.
package value

import (
	"io"
)

var internedSymbols = map[string]*Atom{}

var (
	EOF        = Error("EOF")
	T          = Intern("t")
	NIL        = Intern("nil")                // also: empty list
	unassigned = &Atom{Name: "#[unassigned]"} // uninterned symbol
)

// Value is a runtime representation of Lisp data.
type Value interface {
	Eval(Environment) Value
	Equal(Value) Value
	// If the value has a valid written representation, then this
	// should output an external representation suitable for read.
	Write(io.Writer)
}

func Intern(name string) *Atom {
	v, ok := internedSymbols[name]
	if !ok {
		v = &Atom{Name: name}
		internedSymbols[name] = v
	}
	return v
}

type Atom struct {
	Name string
}

func (v Atom) String() string {
	return v.Name
}

func (v *Atom) Write(w io.Writer) {
	io.WriteString(w, v.Name)
}

func (v *Atom) Eval(env Environment) Value {
	if v == T {
		return T
	}
	if v == NIL {
		return NIL
	}
	if v, ok := env.Lookup(v.Name); ok {
		return v
	}
	return v // auto-quote
}

func (v *Atom) Equal(x Value) Value {
	if _, ok := x.(*Atom); ok {
		if v == x {
			return T
		}
	}
	return NIL
}

// evalList evaluates a proper list of values.
func evalList(env Environment, expr Value) []Value {
	var results []Value
	for next := expr; next != NIL; {
		v, ok := next.(*Cell)
		if !ok {
			Panicf("cannot evaluate an improper list: %s", expr)
		}

		results = append(results, v.Car.Eval(env))

		next = v.Cdr
	}

	return results
}

// evalProgn evaluates a proper list of values, returning the last
// value.
func evalProgn(env Environment, expr Value) Value {
	results := evalList(env, expr)
	return results[len(results)-1]
}

// evalQuote evaluates the quote special form.
func evalQuote(expr *Cell) Value {
	v, ok := expr.Cdr.(*Cell)
	if !ok {
		return Errorf("ill-formed special form: %s", expr)
	}
	if v.Cdr != NIL {
		return Errorf("ill-formed special form: %s", expr)
	}

	return v.Car
}

// evalCond evaluates the cond special form.
func evalCond(env Environment, expr *Cell) Value {
	next := expr
	for next.Cdr != NIL {
		v, ok := next.Cdr.(*Cell)
		if !ok {
			Panicf("ill-formed special form: %s", expr)

		}
		next = v

		clause, ok := next.Car.(*Cell)
		if !ok {
			Panicf("ill-formed special form: %s", expr)
		}

		// if caadr is true, then we want to return the
		// evaluation of the cdadr
		if clause.Car.Eval(env) == T {
			body, ok := clause.Cdr.(*Cell)
			if !ok {
				Panicf("ill-formed clause: %s", clause)
			}
			return evalProgn(env, body)
		}
	}
	return NIL
}

// evalLabel evaluates the label special form.
func evalLabel(env Environment, expr *Cell) Value {
	label, ok := cadr(expr).(*Atom)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	lambda, ok := caddr(expr).(*Cell)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	// Evaluate lambda in an environment where it is able to
	// reference the name defined by the label special form.
	extEnv := NewEnv(env)
	extEnv.Define(label.Name, unassigned)
	fn := makeFunction(extEnv, lambda)

	// Update the binding to the newly created function.
	extEnv.Update(label.Name, fn)

	return fn
}

// evalDefun defines a function permanently.
func evalDefun(env Environment, expr *Cell) *Atom {
	symbol, ok := cadr(expr).(*Atom)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	body, ok := cddr(expr).(*Cell)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	// By defining in the current environment, we add a permanent
	// function, but don't need to find the toplevel. I think this
	// differs from a typical Lisp, but falls within McCarthy's
	// described behaviour.
	env.Define(symbol.Name, unassigned)
	fn := makeFunction(env, Cons(&Atom{Name: "lambda"}, body))
	env.Update(symbol.Name, fn)

	return symbol
}

// makeFunction creates a new function from the lambda special form.
func makeFunction(env Environment, expr *Cell) Function {
	f := &lambdaFunc{}

	// Gather each argument name so we can extend the environment.
	for next := cadr(expr); next != NIL; {
		v, ok := next.(*Cell)
		if !ok {
			Panicf("ill-formed special form: %s", expr)
		}

		atom, ok := v.Car.(*Atom)
		if !ok {
			break
		}
		f.args = append(f.args, atom.Name)

		next = v.Cdr
	}

	// Convert the body into an interpreted function.
	bodyExpr, ok := cddr(expr).(*Cell)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}
	f.fn = func(args []Value) Value {
		if len(args) != len(f.args) {
			Panicf("%s called with %d arguments, but requires %d",
				f, len(args), len(f.args))
		}

		// Arguments are bound in a new environment.
		extEnv := NewEnv(env)
		for i, arg := range f.args {
			extEnv.Define(arg, args[i])
		}

		// Body is evaluated in an implicit progn.
		return evalProgn(extEnv, bodyExpr)
	}

	return f
}
