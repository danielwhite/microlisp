package value

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
	v, ok := arg.(List)
	if !ok || len(v) == 0 {
		Panicf("car: %s is not a pair", arg)
	}
	return v[0]
}

func cdr(arg Value) Value {
	v, ok := arg.(List)
	if !ok {
		Panicf("cdr: %s is not a pair", arg)
	}
	return v[1:]
}

func cons(a Value, b Value) Value {
	switch v := b.(type) {
	case List:
		return append(List{a}, v...)
	default:
		return List{a, b}
	}
}

func list(args []Value) Value {
	if len(args) == 0 {
		return NIL
	}
	return List(append(args, NIL))
}
