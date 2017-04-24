package value

import (
	"io"
)

// Cons constructs an object that holds two values.
//
// Lists can be represented by consing x onto another cons. The y
// value of the last cons must be NIL.
func Cons(x, y Value) *Cell {
	return &Cell{Car: x, Cdr: y}
}

// Cell is an object that holds two values.
type Cell struct {
	Car Value
	Cdr Value
}

// List returns a slice of values in the list represented by a series
// of cons cells. This expects a proper list and will panic if an
// improper list is found.
func (c *Cell) List() []Value {
	var values []Value
	for {
		values = append(values, c.Car)

		if c.Cdr == NIL {
			break
		}

		cdr, ok := c.Cdr.(*Cell)
		if !ok {
			Panicf("expression must be a proper list: %s", c)
		}
		c = cdr
	}
	return values
}

func (c *Cell) Eval(env Environment) Value {
	if c == nil {
		return NIL
	}

	car, ok := c.Car.(*Atom)
	if ok {
		switch car.Name {
		case "quote":
			return evalQuote(c)
		case "cond":
			return evalCond(env, c)
		case "lambda":
			return makeFunction(env, c)
		case "label":
			return evalLabel(env, c)
		case "defun":
			return evalDefun(env, c)
		}
	}

	expr := c.List()
	fn := expr[0].Eval(env)
	args := evalList(env, expr[1:])

	return invoke(fn, args)
}

func (c *Cell) Equal(cmp Value) Value {
	x, ok := cmp.(*Cell)
	if !ok {
		return NIL
	}

	if c.Car.Equal(x.Car) != T {
		return NIL
	}
	if c.Cdr.Equal(x.Cdr) != T {
		return NIL
	}
	return T
}

func (c *Cell) String() string {
	return Sprint(c)
}

func (c *Cell) Write(w io.Writer) {
	io.WriteString(w, "(")
	for {
		c.Car.Write(w)

		// If Cdr is NIL, then finish proper list.
		if c.Cdr == NIL {
			break
		}

		// Check for an improper list.
		cdr, ok := c.Cdr.(*Cell)
		if !ok {
			io.WriteString(w, " . ")
			c.Cdr.Write(w)
			break
		}
		c = cdr // move to next cell

		io.WriteString(w, " ")
	}
	io.WriteString(w, ")")
}
