package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		expr string
		want Node
	}{
		{"foo", &AtomExpr{Name: "foo"}},
		{"()", NIL},
		{"(a)", &ListExpr{Car: &AtomExpr{Name: "a"}, Cdr: NIL}},
		{"(a b)",
			&ListExpr{
				Car: &AtomExpr{Name: "a"},
				Cdr: &ListExpr{
					Car: &AtomExpr{Name: "b"},
					Cdr: NIL}}},
		{"(a (b c))",
			&ListExpr{
				Car: &AtomExpr{Name: "a"},
				Cdr: &ListExpr{
					Car: &ListExpr{
						Car: &AtomExpr{Name: "b"},
						Cdr: &ListExpr{
							Car: &AtomExpr{"c"},
							Cdr: NIL,
						}},
					Cdr: NIL}}},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			got, err := Parse([]byte(tc.expr))
			if err != nil {
				t.Fatalf("parse failed with error: %s", err)
			}

			var buf bytes.Buffer
			Fprint(&buf, got)
			gotExpr := buf.String()
			if tc.expr != gotExpr {
				t.Errorf("want %s, got %s", tc.expr, gotExpr)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want %#v, got %#v", tc.want, got)
			}
		})
	}
}
