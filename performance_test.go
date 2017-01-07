package main

import (
	"bytes"
	"io/ioutil"
	"log"
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

func readFile(name string) []byte {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
