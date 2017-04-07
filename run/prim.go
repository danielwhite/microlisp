package run

import (
	"github.com/danielwhite/microlisp/value"
)

// atom returns T if the value is an atom.
func atom(arg value.Value) value.Value {
	if _, ok := arg.(*value.Atom); ok {
		return value.T
	}
	return value.NIL
}

func equal(a value.Value, b value.Value) value.Value {
	return a.Equal(b)
}

func car(arg value.Value) value.Value {
	v, ok := arg.(value.List)
	if !ok || len(v) == 0 {
		value.Panicf("car: %s is not a pair", arg)
	}
	return v[0]
}

func cdr(arg value.Value) value.Value {
	v, ok := arg.(value.List)
	if !ok {
		value.Panicf("cdr: %s is not a pair", arg)
	}
	return v[1:]
}

func cons(a value.Value, b value.Value) value.Value {
	switch v := b.(type) {
	case value.List:
		return append(value.List{a}, v...)
	default:
		return value.List{a, b}
	}
}

func list(args value.List) value.Value {
	if len(args) == 0 {
		return value.NIL
	}
	return append(args, value.NIL)
}
