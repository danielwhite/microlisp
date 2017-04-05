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
		{`()`, []Token{{LPAREN, ""}, {RPAREN, ""}}},
		{`atom`, []Token{{ATOM, "atom"}}},
		{`(foo bar baz)`, []Token{
			{LPAREN, ""},
			{ATOM, "foo"},
			{ATOM, "bar"},
			{ATOM, "baz"},
			{RPAREN, ""},
		}},
		{`(((a)b)c)`, []Token{
			{LPAREN, ""},
			{LPAREN, ""},
			{LPAREN, ""},
			{ATOM, "a"},
			{RPAREN, ""},
			{ATOM, "b"},
			{RPAREN, ""},
			{ATOM, "c"},
			{RPAREN, ""},
		}},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			got := scanAll(New(strings.NewReader(tc.expr)))
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
