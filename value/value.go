// Package value implements Lisp values and their evaluation.
package value

import (
	"io"
)

var internedSymbols = map[string]*Atom{}

var (
	EOF        = Error("EOF")
	NIL        = Intern("nil")
	T          = Intern("t")
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

type List []Value

func (v List) cadr() Value {
	if len(v) < 2 {
		Panicf("cadr: %s is not a pair", v)
	}
	return v[1]
}

func (v List) String() string {
	return Sprint(v)
}

func (v List) Write(w io.Writer) {
	io.WriteString(w, "(")
	for {
		v[0].Write(w)
		v = v[1:]

		// have we reached the end of the list?
		if len(v) == 1 {
			// is the list improper?
			if v[0] != NIL {
				io.WriteString(w, " . ")
				v[0].Write(w)
			}
			break
		}

		io.WriteString(w, " ")
	}
	io.WriteString(w, ")")
}

func (v List) Eval(env Environment) Value {
	if car, ok := v[0].(*Atom); ok {
		switch car.Name {
		case "quote":
			return evalQuote(env, v)
		case "cond":
			return evalCond(env, v)
		case "lambda":
			return makeFunction(env, v)
		case "label":
			return evalLabel(env, v)
		case "defun":
			return evalDefun(env, v)
		}
	}

	fn := v[0].Eval(env)
	args := v[1:].evalList(env)
	return invoke(fn, args)
}

func (v List) Equal(cmp Value) Value {
	x, ok := cmp.(List)
	if !ok {
		return NIL
	}

	for i := range v {
		if v[i].Equal(x[i]) == NIL {
			return NIL
		}
	}
	return T
}

// evalList evaluates a list of values.
func (v List) evalList(env Environment) []Value {
	if len(v) == 0 {
		return []Value{}
	}

	if v[len(v)-1] != NIL {
		Panicf("evlis: improper argument list")
	}

	results := make([]Value, len(v)-1)
	for i, v := range v[:len(v)-1] {
		results[i] = v.Eval(env)
	}
	return results
}

// evalQuote evaluates the quote special form.
func evalQuote(env Environment, expr List) Value {
	if len(expr) != 3 {
		Panicf("ill-formed special form: %s", Sprint(expr))
	}
	return expr[1]
}

// evalCond evaluates the cond special form.
func evalCond(env Environment, expr List) Value {
	for _, v := range expr[1:] {
		// ensure cadr is a list
		u, ok := v.(List)
		if !ok {
			Panicf("ill-formed special form: %s", expr)
		}

		// if caadr is true, then we want to return the
		// evaluation of the cdadr
		if u[0].Eval(env) == T {
			return u.cadr().Eval(env)
		}
	}

	return NIL
}

// evalLabel evaluates the label special form.
func evalLabel(env Environment, expr List) Value {
	if len(expr) != 4 {
		Panicf("ill-formed special form: %s", expr)
	}

	label, ok := expr[1].(*Atom)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	lambda, ok := expr[2].(List)
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
func evalDefun(env Environment, expr List) *Atom {
	if len(expr) != 5 {
		Panicf("ill-formed special form: %s", expr)
	}

	symbol, ok := expr[1].(*Atom)
	if !ok {
		Panicf("ill-formed special form: %s", expr)
	}

	// By defining in the current environment, we add a permanent
	// function, but don't need to find the toplevel. I think this
	// differs from a typical Lisp, but falls within McCarthy's
	// described behaviour.
	env.Define(symbol.Name, unassigned)
	fn := makeFunction(env, append(List{&Atom{Name: "lambda"}}, expr[2:]...))
	env.Update(symbol.Name, fn)

	return symbol
}

// makeFunction creates a new function from the lambda special form.
func makeFunction(env Environment, expr List) Function {
	if len(expr) < 4 {
		Panicf("ill-formed special form: %s", expr)
	}

	f := &lambdaFunc{}

	if v, ok := expr[1].(List); ok {
		// Each item in the list of arguments must be an atom.
		f.args = make([]string, len(v)-1)
		for i, arg := range v[:len(v)-1] {
			atom, ok := arg.(*Atom)
			if !ok {
				break
			}
			f.args[i] = atom.Name
		}
	} else if expr[1] == NIL {
		// FIXME: This special case is currently necessary
		// because we don't represent lists as cons cells.
		//
		// The only valid case is NIL in which case we have an
		// empty list of args.
		f.args = []string{}
	} else {
		Panicf("ill-formed special form: %s", expr)
	}

	body := expr[2:]
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

		results := body.evalList(extEnv)
		return results[len(results)-1]
	}

	return f
}
