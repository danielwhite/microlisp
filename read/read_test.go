package read_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/run"
	"github.com/danielwhite/microlisp/scan"
	"github.com/danielwhite/microlisp/value"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		expr string
		want value.Value
	}{
		{"foo", &value.Atom{Name: "foo"}},
		{"()", value.NIL},
		{"(a)", value.List{&value.Atom{Name: "a"}, value.NIL}},
		{"(a b)",
			value.List{
				&value.Atom{Name: "a"},
				&value.Atom{Name: "b"},
				value.NIL,
			}},
		{"(a (b c))",
			value.List{
				&value.Atom{Name: "a"},
				value.List{
					&value.Atom{Name: "b"},
					&value.Atom{Name: "c"},
					value.NIL,
				},
				value.NIL,
			}},
		{`; comments
                  (a ; example
                   b)`,
			value.List{
				&value.Atom{Name: "a"},
				&value.Atom{Name: "b"},
				value.NIL,
			}},
		{")", value.Error("unbalanced closed parenthesis")},
		{"(", value.Error("premature EOF")},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			got := run.ReadString(tc.expr)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestReadMultiple(t *testing.T) {
	src := "a (a b) c"
	want := []value.Value{
		&value.Atom{Name: "a"},
		value.List{
			&value.Atom{Name: "a"},
			&value.Atom{Name: "b"},
			value.NIL,
		},
		&value.Atom{Name: "c"},
	}

	got, err := readAll(src)
	if err != nil {
		t.Fatalf("test failed due to read error: %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
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
