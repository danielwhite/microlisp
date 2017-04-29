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
	v, ok := arg.(*Cell)
	if !ok {
		Errorf("car: %s is not a pair", arg)
	}
	return v.Car
}

func cdr(arg Value) Value {
	v, ok := arg.(*Cell)
	if !ok {
		Errorf("cdr: %s is not a pair", arg)
	}
	return v.Cdr
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

// bindings returns an association list mapping each defined symbol in
// an environment to its value.
func bindings(env Environment) Value {
	names := env.Names()
	bindings := make([]Value, len(names))
	for i, name := range names {
		v, _ := env.Lookup(name)
		bindings[i] = Cons(Intern(name), v)
	}
	return list(bindings)
}
