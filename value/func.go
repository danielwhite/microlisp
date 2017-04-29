package value

import (
	"fmt"
)

// invoke applies a list of arguments to a function.
func invoke(v Value, args []Value) Value {
	fn, ok := v.(Function)
	if !ok {
		Panicf("invoke: %s is not a function", v)
	}
	return fn.Invoke(args)
}

type Function interface {
	Value
	Invoke([]Value) Value
}

type lambdaFunc struct {
	args []string
	fn   func([]Value) Value
}

func (f *lambdaFunc) String() string {
	// FIXME: args should be printed with `()`, instead of `[]`.
	return fmt.Sprintf("#[lambda %p %s]", f, f.args)
}

func (f *lambdaFunc) Eval(Environment) Value {
	return f
}

func (f *lambdaFunc) Equal(cmp Value) Value {
	if x, ok := cmp.(*lambdaFunc); ok && f == x {
		return T
	}
	return NIL
}

func (f *lambdaFunc) Invoke(args []Value) Value {
	return f.fn(args)
}

// FuncN creates a Function value from a native Go function that
// accepts a variable number of arguments.
func FuncN(fn func(vs []Value) Value) Function {
	return &nativeFunc{fn: fn}
}

// Func2 creates a Function value from a native Go function that
// accepts a single argument.
func Func1(fn func(Value) Value) Function {
	return FuncN(func(vs []Value) Value {
		if len(vs) != 1 {
			Panicf("called with %d arguments; requires exactly 1 argument", len(vs))
		}
		return fn(vs[0])
	})
}

// Func2 creates a Function value from a native Go function that
// accepts two arguments.
func Func2(fn func(Value, Value) Value) Function {
	return FuncN(func(vs []Value) Value {
		if len(vs) != 2 {
			Panicf("called with %d arguments; requires exactly 1 argument", len(vs))
		}
		return fn(vs[0], vs[1])
	})
}

// nativeFunc holds a native function that accepts a variable number
// arugments.
type nativeFunc struct {
	fn func([]Value) Value
}

func (f *nativeFunc) String() string {
	return fmt.Sprintf("#[compiled-function %p]", f)
}

func (f *nativeFunc) Eval(Environment) Value {
	return f
}

func (f *nativeFunc) Invoke(args []Value) Value {
	return f.fn(args)
}

func (f *nativeFunc) Equal(cmp Value) Value {
	if x, ok := cmp.(*nativeFunc); ok && f == x {
		return T
	}
	return NIL
}
