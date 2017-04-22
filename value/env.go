package value

import "errors"

// An environment maintains a set of bindings that are typically
// referenced when evaluating a Lisp expression.
type Environment interface {
	// Lookup the value of a symbol.
	Lookup(name string) (Value, bool)
	// Define a new symbol.
	Define(name string, value Value)
	// Update the value of a symbol.
	Update(name string, value Value) error
}

var equalFn = Func2(equal)

var DefaultEnvironment = &env{
	env: map[string]Value{
		"t":     T,
		"nil":   NIL,
		"atom":  Func1(atom),
		"eq":    equalFn, // alias
		"equal": equalFn,
		"car":   Func1(car),
		"cdr":   Func1(cdr),
		"cons":  Func2(cons),
		"list":  FuncN(list),
	},
}

// NewEnv returns a new environment that extends the bindings of
// parent. If parent is nil, then this is a toplevel environment.
func NewEnv(parent Environment) Environment {
	return &env{
		env:    make(map[string]Value),
		parent: parent,
	}
}

type env struct {
	env    map[string]Value
	parent Environment
}

// Define implements the Environment interface.
func (e *env) Define(name string, value Value) {
	e.env[name] = value
}

// Lookup implements the Environment interface.
func (e *env) Lookup(name string) (Value, bool) {
	if v, ok := e.env[name]; ok {
		return v, true
	}

	if e.parent == nil {
		return nil, false
	}
	return e.parent.Lookup(name)
}

// Update implements the Environment interface.
func (e *env) Update(name string, value Value) error {
	if _, ok := e.env[name]; ok {
		e.env[name] = value
		return nil
	}

	// Only define can introduce new bindings.
	if e.parent == nil {
		return errors.New("no such binding")
	}

	return e.parent.Update(name, value)
}
