package main

import (
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	type pair struct {
		tok Token
		lit string
	}
	testCases := []struct {
		expr string
		want []pair
	}{
		{`()`, []pair{{LPAREN, ""}, {RPAREN, ""}}},
		{`atom`, []pair{{ATOM, "atom"}}},
		{`(foo bar baz)`, []pair{
			{LPAREN, ""},
			{ATOM, "foo"},
			{ATOM, "bar"},
			{ATOM, "baz"},
			{RPAREN, ""},
		}},
		{`(((a)b)c)`, []pair{
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
			s := NewScanner(strings.NewReader(tc.expr))

			var got []pair
			for {
				tok, lit := s.Scan()
				if tok == EOF {
					break
				}
				got = append(got, pair{tok, lit})
			}
			for i, pair := range got {
				wantTok := tc.want[i].tok
				wantLit := tc.want[i].lit
				if pair.tok != wantTok || pair.lit != wantLit {
					t.Errorf("want (%s,%q), got (%s,%q)",
						wantTok, wantLit,
						pair.tok, pair.lit)
				}
			}
		})
	}
}
