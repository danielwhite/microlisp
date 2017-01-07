package main

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		expr string
		want Node
	}{
		{"foo", &AtomExpr{Name: "foo"}},
		{"()", NIL},
		{"(a)", &ListExpr{&AtomExpr{Name: "a"}, NIL}},
		{"(a b)",
			&ListExpr{
				&AtomExpr{Name: "a"},
				&ListExpr{
					&AtomExpr{Name: "b"},
					NIL}}},
		{"(a (b c))",
			&ListExpr{
				&AtomExpr{Name: "a"},
				&ListExpr{
					&ListExpr{
						&AtomExpr{Name: "b"},
						&ListExpr{
							&AtomExpr{"c"},
							NIL,
						}},
					NIL}}},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			s := NewScanner(strings.NewReader(tc.expr))

			got, err := Read(s)
			if err != nil {
				t.Fatalf("test failed due to read error: %s", err)
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

func TestReadMultiple(t *testing.T) {
	src := "a (a b) c"
	want := []Node{
		&AtomExpr{Name: "a"},
		&ListExpr{
			&AtomExpr{Name: "a"},
			&ListExpr{&AtomExpr{Name: "b"}, NIL}},
		&AtomExpr{Name: "c"},
	}

	s := NewScanner(strings.NewReader(src))
	got, err := readAll(s)
	if err != nil {
		t.Fatalf("test failed due to read error: %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v", want, got)
	}
}

func readAll(s *Scanner) ([]Node, error) {
	var nodes []Node
	for {
		node, err := Read(s)

		if err == io.EOF {
			return nodes, nil
		} else if err != nil {
			return nodes, err
		}

		nodes = append(nodes, node)
	}
}

func TestReadError(t *testing.T) {
	testCases := []struct {
		expr    string
		wantErr string
	}{
		{")", "unbalanced closed parenthesis"},
		{"(", "premature EOF"},
	}
	for _, tc := range testCases {
		s := NewScanner(strings.NewReader(tc.expr))
		_, err := Read(s)
		if err == nil {
			t.Fatalf("no error, expected: %s", err)
		}
		if tc.wantErr != err.Error() {
			t.Fatalf("want %q, got %v", tc.wantErr, err)
		}
	}
}
