package read_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/run"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

func TestRead(t *testing.T) {
	a, b, c := value.Intern("a"), value.Intern("b"), value.Intern("c")

	testCases := []struct {
		expr string
		want value.Value
	}{
		{"a", a},
		{"()", value.NIL},
		{"(a)", value.Cons(a, value.NIL)},
		{"(a b)", value.Cons(a, value.Cons(b, value.NIL))},
		{"(a b c)", value.Cons(a, value.Cons(b, value.Cons(c, value.NIL)))},
		{"(a (b c))",
			value.Cons(a,
				value.Cons(
					value.Cons(b, value.Cons(c, value.NIL)),
					value.NIL))},
		{`; comments
                  (a ; example
                   b)`,
			value.Cons(a, value.Cons(b, value.NIL))},
		{")", value.Error("unbalanced closed parenthesis")},
		{"(", value.Error("premature EOF")},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			got := run.ReadString(tc.expr)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want %s, got %s", tc.want, got)
			}
		})
	}
}

func TestReadMultiple(t *testing.T) {
	a, b, c := value.Intern("a"), value.Intern("b"), value.Intern("c")

	src := "a (a b) c"
	want := []value.Value{
		a,
		value.Cons(a, value.Cons(b, value.NIL)),
		c,
	}

	got, err := readAll(src)
	if err != nil {
		t.Fatalf("test failed due to read error: %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %s, got %s", want, got)
	}
}

func readAll(text string) (values []value.Value, err error) {
	scanner := scan.New(strings.NewReader(text))
	reader := read.New(scanner)
	for {
		v := reader.Read()

		if v == value.EOF {
			return
		}
		if v, ok := v.(value.Error); ok {
			err = errors.New(v.Error())
			return
		}

		values = append(values, v)
	}
}
