package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/run"
	"github.com/danielwhite/microlisp/scan"
	"github.com/danielwhite/microlisp/value"
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
	v.Write(&buf)
	b.SetBytes(int64(buf.Len()))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		v.Write(&buf)
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
		{"cons", "1 2"},
		{"equal", "(quote (a (b c) d (e (f) g))) (quote (a (b c) d (e (f) g)))"},
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
	b, err := ioutil.ReadFile(name)
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
