package value

func NewEnv(e map[string]Value) Environment {
	return &env{env: e}
}

type env struct {
	env  map[string]Value
	next Environment
}

// extend returns an extended environment.
func (e *env) Extend(extEnv map[string]Value) Environment {
	return &env{
		env:  extEnv,
		next: e,
	}
}

// Lookup implements the Environment interface.
func (e *env) Lookup(name string) (Value, bool) {
	if v, ok := e.env[name]; ok {
		return v, true
	}

	if e.next == nil {
		return nil, false
	}
	return e.next.Lookup(name)
}
