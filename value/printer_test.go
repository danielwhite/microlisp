package value

import (
	"testing"
)

func TestSprint(t *testing.T) {
	testCases := []struct {
		value Value
		want  string
	}{
		{Intern("a"), "a"},
		{NIL, "nil"},
		{List{NIL, NIL}, "(nil)"},
		{List{Intern("a"), NIL}, "(a)"},
		{List{
			Intern("a"),
			Intern("b"),
			NIL},
			"(a b)"},
		// Printing of improper lists.
		{List{
			Intern("a"),
			Intern("b")},
			"(a . b)"},
		{List{
			Intern("a"),
			Intern("b"),
			Intern("c")},
			"(a b . c)"},
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
