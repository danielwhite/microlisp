package value

import (
	"testing"
)

func TestSprint(t *testing.T) {
	a, b, c := Intern("a"), Intern("b"), Intern("c")

	testCases := []struct {
		value Value
		want  string
	}{
		{a, "a"},
		{NIL, "nil"},
		{Cons(NIL, NIL), "(nil)"},
		{Cons(a, NIL), "(a)"},
		{Cons(a, Cons(b, NIL)), "(a b)"},
		// Printing of improper lists.
		{Cons(a, b), "(a . b)"},
		{Cons(a, Cons(b, c)), "(a b . c)"},
	}
	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			got := Sprint(tc.value)
			if tc.want != got {
				t.Errorf("want:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}
}
