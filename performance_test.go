package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var src = readFile("testdata/eval.lisp")

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		s := NewScanner(bytes.NewReader(src))
		if _, err := Read(s); err != nil {
			b.Fatalf("benchmark failed due to parse error: %s", err)
		}
	}
}

func BenchmarkFprint(b *testing.B) {
	s := NewScanner(bytes.NewReader(src))
	node, err := Read(s)
	if err != nil {
		b.Fatalf("benchmark failed due to parse error: %s", err)
	}

	// Initial print to allocate underlying buffer.
	var buf bytes.Buffer
	Fprint(&buf, node)
	b.SetBytes(int64(buf.Len()))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		Fprint(&buf, node)
	}
}

func BenchmarkEval(b *testing.B) {
	s := NewScanner(bytes.NewReader(src))
	node, err := Read(s)
	if err != nil {
		b.Fatalf("benchmark failed due to parse error: %s", err)
	}

	for i := 0; i < b.N; i++ {
		Eval(node)
	}
}

func BenchmarkInvoke(b *testing.B) {
	testCases := []struct {
		name string
		args string
	}{
		{"list", "1 2 3"},
		{"car", "(quote (x y z))"},
		{"cdr", "(quote (a b c))"},
		{"cons", "1 2"},
		{"equal", "(quote (a (b c) d (e (f) g))) (quote (a (b c) d (e (f) g)))"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			s := NewScanner(strings.NewReader(fmt.Sprintf("(%s %s)", tc.name, tc.args)))
			node, err := Read(s)
			if err != nil {
				b.Fatalf("benchmark failed due to parse error: %s", err)
			}

			for i := 0; i < b.N; i++ {
				Eval(node)
			}
		})
	}
}

func readFile(name string) []byte {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
