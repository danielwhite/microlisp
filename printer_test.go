package main

import (
	"bytes"
	"testing"
)

func TestFprint(t *testing.T) {
	testCases := []struct {
		node Node
		want string
	}{
		{&AtomExpr{Name: "a"}, "a"},
		{NIL, "()"},
		{&ListExpr{Car: NIL, Cdr: NIL}, "(())"},
		{&ListExpr{Car: &AtomExpr{Name: "a"}, Cdr: NIL}, "(a)"},
		{&ListExpr{
			Car: &AtomExpr{Name: "a"},
			Cdr: &ListExpr{Car: &AtomExpr{Name: "b"}, Cdr: NIL}},
			"(a b)"},
		// Printing of improper lists.
		{&ListExpr{Car: &AtomExpr{Name: "a"}, Cdr: &AtomExpr{Name: "b"}}, "(a . b)"},
		{&ListExpr{
			Car: &AtomExpr{Name: "a"},
			Cdr: &ListExpr{
				Car: &AtomExpr{Name: "b"},
				Cdr: &AtomExpr{Name: "c"}}},
			"(a b . c)"},
	}
	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			var buf bytes.Buffer
			if err := Fprint(&buf, tc.node); err != nil {
				t.Fatalf("printing node failed: %s", err)
			}
			got := buf.String()
			if tc.want != got {
				t.Errorf("want:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}
}
