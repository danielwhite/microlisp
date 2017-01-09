package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	testCases := []struct {
		expr    string // expression to evalute
		want    string // value of the expression
		wantErr string // error of a bad expression
	}{
		{"(quote a)", "a", ""},
		{"(quote (a b c))", "(a b c)", ""},
		{"(quote)", "", "ill-formed special form: (quote)"},
		{"(quote a b)", "", "ill-formed special form: (quote a b)"},
		{"(car (quote (1 2)))", "1", ""},
		{"(car (quote 1))", "", "car: 1 is not a pair"},
		{"(cdr (quote (1 2)))", "(2)", ""},
		{"(cdr (quote 1))", "", "cdr: 1 is not a pair"},
		{"(list)", "()", ""},
		{"(list 1)", "(1)", ""},
		{"(list 1 2)", "(1 2)", ""},
		{"(list 1 2 3)", "(1 2 3)", ""},
		{"((1 2) 3 4)", "", "invoke: 1 is not a function"},
		{"((car (list cdr car)) (quote (1 2 3)))", "(2 3)", ""},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			s := NewScanner(strings.NewReader(tc.expr))
			expr, err := Read(s)
			if err != nil {
				t.Fatalf("test failed due to read error: %s", err)
			}

			// Evaluate the expression.
			value, err := Eval(expr)
			if tc.wantErr == "" && err != nil {
				t.Fatalf("test failed due to eval error: %s", err)
			}
			if tc.wantErr != "" && (err == nil || tc.wantErr != err.Error()) {
				t.Fatalf("want %q, got %v", tc.wantErr, err)
			}

			var buf bytes.Buffer
			Fprint(&buf, value)
			got := buf.String()

			if tc.want != got {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}
