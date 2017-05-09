package value

import (
	"fmt"
)

// invoke applies a list of arguments to a function.
func invoke(v Value, args []Value) Value {
	fn, ok := v.(Function)
	if !ok {
		Errorf("invoke: %s is not a function", v)
	}
	return fn.Invoke(args)
}

type Function interface {
	Value
	Invoke([]Value) Value
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
		assertArgs(1, len(vs))
		return fn(vs[0])
	})
}

// Func2 creates a Function value from a native Go function that
// accepts two arguments.
func Func2(fn func(Value, Value) Value) Function {
	return FuncN(func(vs []Value) Value {
		assertArgs(2, len(vs))
		return fn(vs[0], vs[1])
	})
}

// FuncX creates a Function value from a native Go function that
// accepts a specified number of arguments.
func FuncX(n int, fn func([]Value) Value) Function {
	return FuncN(func(vs []Value) Value {
		assertArgs(n, len(vs))
		return fn(vs)
	})
}

func assertArgs(want, got int) {
	if want == got {
		return
	}
	switch {
	case want == 1:
		Errorf("called with %d arguments; requires exactly 1 argument", got)
	case got == 1:
		Errorf("called with 1 arguments; requires exactly %d arguments", want)
	default:
		Errorf("called with %d arguments; requires exactly %d arguments", got, want)
	}
}

// nativeFunc holds a native function that accepts a variable number
// arugments.
type nativeFunc struct {
	fn func([]Value) Value
}

func (f *nativeFunc) String() string {
	return fmt.Sprintf("#[compiled-function %p]", f)
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
