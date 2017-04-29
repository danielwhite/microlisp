package value

import (
	"errors"
	"fmt"
)

// An environment maintains a set of bindings that are typically
// referenced when evaluating a Lisp expression.
type Environment interface {
	Value
	// Returns the list of defined symbols.
	Names() []string
	// Lookup the value of a symbol.
	Lookup(name string) (Value, bool)
	// Define a new symbol.
	Define(name string, value Value)
	// Update the value of a symbol.
	Update(name string, value Value) error
}

var equalFn = Func2(equal)

// SystemEnvironment is the toplevel environment where primitives are
// defined.
var SystemEnvironment = &env{
	env: map[string]Value{
		"atom":   Func1(atom),
		"null":   Func1(null),
		"eq":     equalFn, // alias
		"equal":  equalFn,
		"car":    Func1(car),
		"cdr":    Func1(cdr),
		"caar":   Func1(caar),
		"cadr":   Func1(cadr),
		"cddr":   Func1(cddr),
		"caddr":  Func1(caddr),
		"cadar":  Func1(cadar),
		"caddar": Func1(caddar),
		"cons":   Func2(func(x, y Value) Value { return Cons(x, y) }),
		"list":   FuncN(list),

		// Environment Primitives
		"environment-bindings": EnvFunc(bindings),
	},
}

// FuncEnv creates a Function value from a native Go function that
// only accepts environments.
func EnvFunc(fn func(Environment) Value) Function {
	return Func1(func(v Value) Value {
		if env, ok := v.(Environment); ok {
			return fn(env)
		}
		return Errorf("%s is not an environment", v)
	})
}

func init() {
	SystemEnvironment.Define("system-environment", SystemEnvironment)
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

func (e *env) String() string {
	return fmt.Sprintf("#[env %p %d]", e, len(e.env))
}

// Equal implments the Value interface, and returns T for the same
// environment.
func (e *env) Equal(cmp Value) Value {
	if x, ok := cmp.(*env); ok && e == x {
		return T
	}
	return NIL
}

// Eval implements the Value interface.
func (e *env) Eval(Environment) Value {
	return e
}

// Names implements the Environment interface, returning a list of all
// defined symbols.
func (e *env) Names() []string {
	names := make([]string, 0, len(e.env))
	for k := range e.env {
		names = append(names, k)
	}
	return names
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
