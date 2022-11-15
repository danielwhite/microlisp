package value

import (
	"fmt"
	"strings"
)

// atom returns T if the value is an atom.
func atom(arg Value) Value {
	if _, ok := arg.(*Atom); ok {
		return T
	}
	return NIL
}

// null returns T if the value is NIL
func null(v Value) Value {
	if v == NIL {
		return T
	}
	return NIL
}

func equal(a Value, b Value) Value {
	return a.Equal(b)
}

func car(arg Value) Value {
	v, ok := arg.(*Cell)
	if !ok {
		Errorf("car: %s is not a pair", arg)
	}
	return v.Car
}

func cdr(arg Value) Value {
	v, ok := arg.(*Cell)
	if !ok {
		Errorf("cdr: %s is not a pair", arg)
	}
	return v.Cdr
}

func caar(v Value) Value {
	return car(car(v))
}

func cadr(v Value) Value {
	return car(cdr(v))
}

func cddr(v Value) Value {
	return cdr(cdr(v))
}

func caddr(v Value) Value {
	return car(cdr(cdr(v)))
}

func cadar(v Value) Value {
	return car(cdr(car(v)))
}

func caddar(v Value) Value {
	return car(cdr(cdr(car(v))))
}

func list(args []Value) Value {
	if len(args) == 0 {
		return NIL
	}

	var last Value = NIL
	for i := len(args) - 1; i >= 0; i-- {
		last = Cons(args[i], last)
	}
	return last
}

// apply a list as arguments to a function.
//
// This might be simpler by implementing it in terms of invoke.
func apply(vs []Value) Value {
	if len(vs) < 1 {
		Errorf("called with %d arguments; requires at least 1 argument", len(vs))
	}

	fn, rest := vs[0], vs[1:]
	if len(rest) == 0 {
		return invoke(vs[0], []Value{})
	}

	// Each initial argument is prepended onto the final cons cell.
	last := rest[len(rest)-1]
	head := last
	for i := len(vs) - 2; i > 0; i-- {
		head = Cons(vs[i], head)
	}

	// If the final cell is not a list, report error with all the arguments.
	if _, ok := last.(*Cell); !ok {
		Errorf("apply: improper argument list: %s", head)
	}

	var args []Value
	head.(*Cell).Walk(func(v Value) {
		args = append(args, v)
	})
	return invoke(fn, args)
}

// bindings returns an association list mapping each defined symbol in
// an environment to its value.
func bindings(env Environment) Value {
	names := env.Names()
	bindings := make([]Value, len(names))
	for i, name := range names {
		v, _ := env.Lookup(name)
		bindings[i] = Cons(Intern(name), v)
	}
	return list(bindings)
}

// raiseError pretty-prints the values passed, and throws a
// recoverable error.
func raiseError(vs []Value) Value {
	Errorf(strings.Trim(fmt.Sprintf("%s", vs), "[]"))
	panic("not possible")
}

// trapError returns the value of an invoked function. If an error is
// raised, the error value is instead returned.
func trapError(fn Value) (v Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(Error)
		}
	}()
	v = invoke(fn, []Value{})
	return
}
