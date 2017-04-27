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
	switch v := arg.(type) {
	case *Cell:
		return v.Car
	case Error:
		return arg
	default:
		return Errorf("car: %s is not a pair", arg)
	}
}

func cdr(arg Value) Value {
	switch v := arg.(type) {
	case *Cell:
		return v.Cdr
	case Error:
		return arg
	default:
		return Errorf("cdr: %s is not a pair", arg)
	}
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
