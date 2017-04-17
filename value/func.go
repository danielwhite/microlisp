package value

import (
	"fmt"
	"io"
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

func (f *lambdaFunc) Write(w io.Writer) {
	// FIXME: args should be printed with `()`, instead of `[]`.
	fmt.Fprintf(w, "#[lambda %p %s]", f, f.args)
}

func (f *lambdaFunc) String() string {
	return Sprint(f)
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

func Func1(fn func(Value) Value) Function {
	return &func1{fn}
}

// func1 represents a function that accepts a single argument.
type func1 struct {
	fn func(Value) Value
}

func (f *func1) Write(w io.Writer) {
	fmt.Fprintf(w, "#[compiled-function %p]", f)
}

func (f *func1) String() string {
	return Sprint(f)
}

func (f *func1) Eval(Environment) Value {
	return f
}

func (f *func1) Invoke(args []Value) Value {
	if len(args) != 1 {
		Panicf("called with %d arguments; requires exactly 1 argument", len(args))
	}
	return f.fn(args[0])
}

func (f *func1) Equal(cmp Value) Value {
	if x, ok := cmp.(*func1); ok && f == x {
		return T
	}
	return NIL
}

func Func2(fn func(Value, Value) Value) Function {
	return &func2{fn}
}

// func2 represents a function that accepts exactly two arguments.
type func2 struct {
	fn func(Value, Value) Value
}

func (f *func2) Write(w io.Writer) {
	fmt.Fprintf(w, "#[compiled-function %p]", f)
}

func (f *func2) String() string {
	return Sprint(f)
}

func (f *func2) Eval(Environment) Value {
	return f
}

func (f *func2) Invoke(args []Value) Value {
	if len(args) != 2 {
		Panicf("called with %d arguments; requires exactly 1 argument", len(args))
	}
	return f.fn(args[0], args[1])
}

func (f *func2) Equal(cmp Value) Value {
	if x, ok := cmp.(*func2); ok && f == x {
		return T
	}
	return NIL
}

func FuncN(fn func([]Value) Value) Function {
	return &funcN{fn}
}

// funcN represents a function that accepts a variable number
// arugments.
type funcN struct {
	fn func([]Value) Value
}

func (f *funcN) Write(w io.Writer) {
	fmt.Fprintf(w, "#[compiled-function %p]", f)
}

func (f *funcN) String() string {
	return Sprint(f)
}

func (f *funcN) Eval(Environment) Value {
	return f
}

func (f *funcN) Invoke(args []Value) Value {
	return f.fn(args)
}

func (f *funcN) Equal(cmp Value) Value {
	if x, ok := cmp.(*funcN); ok && f == x {
		return T
	}
	return NIL
}
