package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/run"
	"github.com/danielwhite/microlisp/scan"
)

type Matcher interface {
	MatchString(string) bool
}

type stringMatcher string

func (m stringMatcher) MatchString(s string) bool {
	return string(m) == s
}

func TestEval(t *testing.T) {
	testCases := []struct {
		expr string      // expression to evalute
		want interface{} // value of the expression; string or Matcher
	}{
		{"t", "t"},
		{"nil", "nil"},

		{"(equal t (quote t))", "t"},

		{"(atom t)", "t"},
		{"(atom nil)", "t"},
		{"(atom ())", "t"},
		{"(atom (cons 1 2))", "nil"},
		{"(atom (cons 1 (cons 2 nil)))", "nil"},

		{"(equal (car (quote (a b))) (quote a))", "t"},
		{"(equal (cdr (quote (a b))) (quote a))", "nil"},
		{"(equal (quote (a (b c) d)) (list (quote a) (quote (b c)) (quote d)))", "t"},
		{"(equal car car)", "t"},
		{"(equal car cdr)", "nil"},

		{"(quote a)", "a"},
		{"(quote (a b c))", "(a b c)"},
		{"(quote)", "#[error: ill-formed special form: (quote)]"},
		{"(quote a b)", "#[error: ill-formed special form: (quote a b)]"},

		{"(car (quote (1 2)))", "1"},
		{"(car (quote 1))", "#[error: car: 1 is not a pair]"},

		{"(cdr (quote (1 2)))", "(2)"},
		{"(cdr (quote 1))", "#[error: cdr: 1 is not a pair]"},

		{"(cons 1 2)", "(1 . 2)"},
		{"(cons 1 (cons 2 ()))", "(1 2)"},
		{"(cons (quote a) (quote (b c)))", "(a b c)"},

		{"(list)", "nil"},
		{"(list 1)", "(1)"},
		{"(list 1 2)", "(1 2)"},
		{"(list 1 2 3)", "(1 2 3)"},

		{"((1 2) 3 4)", "#[error: invoke: 1 is not a function]"},
		{"((car (list cdr car)) (quote (1 2 3)))", "(2 3)"},

		{"(cond ((atom (quote a)) (quote b)) ((quote t) (quote c)))", "b"},
		{"(cond ((atom car) (quote b)) ((quote t) (quote c)))", "c"},

		{"((lambda () (quote a)))", "a"},
		{"((lambda (a) a) b)", "b"},
		{"((lambda (x y) (cons (car x) y)) (quote (a b)) (cdr (quote (c d))))",
			"(a d)"},

		// Ensure inner Lambda access extended environment.
		{"(((lambda (x) (lambda (y) (cons x y))) a) b)", "(a . b)"},
		// Ensure body evaluation is equivalent to `progn`.
		{"((lambda (x) a x c) c)", "c"},

		{"(lambda)", "#[error: ill-formed special form: (lambda)]"},
		{"(lambda ())", "#[error: ill-formed special form: (lambda nil)]"},
		{"((lambda (x) x) 1 2)", regexp.MustCompile(`#\[error: #\[lambda .* \[x\]\] called with 2 arguments, but requires 1\]`)},

		{`((label ff (lambda (x)
                               (cond ((atom x) x)
                                     ((quote t) (ff (car x))))))
                   (quote ((a b) c)))`,
			"a"},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			// Ensure that test case provides a matcher.
			var matcher Matcher
			switch v := tc.want.(type) {
			case string:
				matcher = stringMatcher(v)
			case Matcher:
				matcher = v
			default:
				t.Fatalf("test requires a matcher, got %#v", tc.want)
			}

			// Read and evaluate an expression.
			scanner := scan.New(strings.NewReader(tc.expr))
			reader := read.New(scanner)
			expr := reader.Read()
			value := run.Eval(expr)

			// Print result of expression.
			var buf bytes.Buffer
			value.Write(&buf)
			got := buf.String()

			if !matcher.MatchString(got) {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}
