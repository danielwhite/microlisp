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
		{"(quote)", "", "cadr: <nil> is not a pair"},
		{"(quote a b)", "", "ill-formed special form"},
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
