package value

import (
	"bytes"
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

// Walk traverses calls fn(Car) for each cell in a list.
func (c *Cell) Walk(fn func(Value)) {
	cur := c
	for {
		fn(cur.Car)

		// The last cell has been reached.
		if cur.Cdr == NIL {
			return
		}

		next, ok := cur.Cdr.(*Cell)
		if !ok {
			Errorf("cannot evaluate an improper list: %s", c)
		}
		cur = next
	}
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
	var buf bytes.Buffer

	buf.WriteByte('(')
	for {
		buf.WriteString(c.Car.String())

		// If Cdr is NIL, then finish proper list.
		if c.Cdr == NIL {
			break
		}

		// Check for an improper list.
		cdr, ok := c.Cdr.(*Cell)
		if !ok {
			buf.WriteString(" . ")
			buf.WriteString(c.Cdr.String())
			break
		}
		c = cdr // move to next cell

		buf.WriteByte(' ')
	}
	buf.WriteByte(')')

	return buf.String()
}
