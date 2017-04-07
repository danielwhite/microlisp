package run

import (
	"github.com/danielwhite/microlisp/value"
)

var DefaultEnvironment = mapEnv{
	"t":     value.T,
	"nil":   value.NIL,
	"atom":  value.Func1(atom),
	"equal": value.Func2(equal),
	"car":   value.Func1(car),
	"cdr":   value.Func1(cdr),
	"cons":  value.Func2(cons),
	"list":  value.FuncN(list),
}

// Eval applies rules to an expression, and returns an expression that
// is the value.
func Eval(expr value.Value) (v value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(value.Error)
		}
	}()
	v = expr.Eval(DefaultEnvironment)
	return
}

type mapEnv map[string]value.Value

// Lookup implements the value.Environment interface.
func (e mapEnv) Lookup(name string) (value.Value, bool) {
	v, ok := e[name]
	return v, ok
}
