package value

import "testing"

func TestApply(t *testing.T) {
	LIST := FuncN(list)
	A, B, C, D := Intern("A"), Intern("B"), Intern("C"), Intern("D")

	testCases := []struct {
		Args []Value
		Want Value
	}{
		{[]Value{LIST}, NIL},
		{[]Value{LIST, list([]Value{A})}, list([]Value{A})},
		{[]Value{LIST, A, list([]Value{B})}, list([]Value{A, B})},

		{[]Value{LIST, A}, Error("apply: improper argument list: A")},
		{[]Value{LIST, A, B}, Error("apply: improper argument list: (A . B)")},
		{[]Value{LIST, A, B, C}, Error("apply: improper argument list: (A B . C)")},
		{[]Value{LIST, A, B, C, D}, Error("apply: improper argument list: (A B C . D)")},
	}
	for _, tc := range testCases {
		got := trapError(FuncN(func(_ []Value) Value {
			return apply(tc.Args)
		}))

		if got.Equal(tc.Want) != T {
			t.Errorf("Apply(%+v) = %v, want %v", tc.Args, got, tc.Want)
		}
	}
}
