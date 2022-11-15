package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/run"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

var src = readFile("testdata/eval.lisp")

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		scanner := scan.New(bytes.NewReader(src))
		reader := read.New(scanner)
		v := reader.Read()
		if err, ok := v.(value.Error); ok {
			b.Fatalf("benchmark failed due to parse error: %s", err)
		}
	}
}

func BenchmarkFprint(b *testing.B) {
	scanner := scan.New(bytes.NewReader(src))
	reader := read.New(scanner)
	v := reader.Read()
	if err, ok := v.(value.Error); ok {
		b.Fatalf("benchmark failed due to parse error: %s", err)
	}

	// Initial print to allocate underlying buffer.
	var buf bytes.Buffer
	buf.WriteString(v.String())
	b.SetBytes(int64(buf.Len()))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.WriteString(v.String())
	}
}

func BenchmarkEval(b *testing.B) {
	scanner := scan.New(bytes.NewReader(src))
	reader := read.New(scanner)
	v := reader.Read()
	if err, ok := v.(value.Error); ok {
		b.Fatalf("benchmark failed due to parse error: %s", err)
	}

	for i := 0; i < b.N; i++ {
		run.Eval(v)
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
		{"caar", "(quote ((a . 1) (b . 2) (c . 3)))"},
		{"cadr", "(quote (a b c))"},
		{"cddr", "(quote (a b c))"},
		{"caddr", "(quote (a b c))"},
		{"cadar", "(quote ((a a') b c))"},
		{"caddar", "(quote ((a a' a'') b c))"},
		{"cons", "1 2"},
		{"equal", "(quote (a (b c) d (e (f) g))) (quote (a (b c) d (e (f) g)))"},
		{"lambda", "(a) a"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			src := fmt.Sprintf("(%s %s)", tc.name, tc.args)
			scanner := scan.New(strings.NewReader(src))
			reader := read.New(scanner)
			v := reader.Read()
			if err, ok := v.(value.Error); ok {
				b.Fatalf("benchmark failed due to parse error: %s", err)
			}

			for i := 0; i < b.N; i++ {
				if _, err := unwrapEval(v); err != nil {
					b.Fatalf("benchmark failed due to eval error: %s", err)
				}
			}
		})
	}
}

func readFile(name string) []byte {
	b, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func unwrapEval(expr value.Value) (value.Value, error) {
	v := run.Eval(expr)
	if err, ok := v.(value.Error); ok {
		return nil, err
	}
	return v, nil
}
