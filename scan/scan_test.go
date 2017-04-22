package scan

import (
	"reflect"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	testCases := []struct {
		expr string
		want []Token
	}{
		{`()`, []Token{{LeftParen, ""}, {RightParen, ""}}},
		{`atom`, []Token{{Atom, "atom"}}},
		{`(foo bar baz)`, []Token{
			{LeftParen, ""},
			{Atom, "foo"},
			{Atom, "bar"},
			{Atom, "baz"},
			{RightParen, ""},
		}},
		{`(((a)b)c)`, []Token{
			{LeftParen, ""},
			{LeftParen, ""},
			{LeftParen, ""},
			{Atom, "a"},
			{RightParen, ""},
			{Atom, "b"},
			{RightParen, ""},
			{Atom, "c"},
			{RightParen, ""},
		}},
		{`(list ;; comment
                    ;; some values
                    a
                    b)`, []Token{
			{LeftParen, ""},
			{Atom, "list"},
			{Comment, ";; comment"},
			{Comment, ";; some values"},
			{Atom, "a"},
			{Atom, "b"},
			{RightParen, ""},
		}},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			scanner := New(strings.NewReader(tc.expr))

			got := scanAll(scanner)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func scanAll(s *Scanner) []Token {
	var toks []Token
	for {
		tok := s.Next()
		if tok.Type == EOF {
			break
		}
		toks = append(toks, tok)
	}
	return toks
}
