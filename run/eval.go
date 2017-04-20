package run

import (
	"github.com/danielwhite/microlisp/value"
)

// Eval applies rules to an expression, and returns an expression that
// is the value.
func Eval(expr value.Value) (v value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(value.Error)
		}
	}()
	v = expr.Eval(value.DefaultEnvironment)
	return
}
