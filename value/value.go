// Package value implements Lisp values and their evaluation.
package value

var internedSymbols = map[string]*Atom{}

var (
	EOF = Error("EOF")
	T   = Intern("t")
	NIL = Intern("nil") // also: empty list
)

// Value is a runtime representation of Lisp data.
type Value interface {
	Equal(Value) Value
	// If the value has a valid written representation, then this
	// should output an external representation suitable for read.
	String() string
}

func Intern(name string) *Atom {
	v, ok := internedSymbols[name]
	if !ok {
		v = &Atom{Name: name}
		internedSymbols[name] = v
	}
	return v
}

type Atom struct {
	Name string
}

func (v Atom) String() string {
	return v.Name
}

func (v *Atom) Equal(x Value) Value {
	if _, ok := x.(*Atom); ok {
		if v == x {
			return T
		}
	}
	return NIL
}
