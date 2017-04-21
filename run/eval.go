package run

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/scan"
	"github.com/danielwhite/microlisp/value"
)

// Eval applies rules to an expression, and returns an expression that
// is the value.
func Eval(expr value.Value) (v value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(value.Error)
		}
	}()
	v = expr.Eval(value.DefaultEnvironment)
	return
}

// Run a REPL loop, reading expressions from a reader, and writing the
// evaluated values to a writer.
func Run(r io.Reader, w io.Writer) error {
	scanner := scan.New(bufio.NewReader(r))
	reader := read.New(scanner)
	for {
		// Read the next expression from the input.
		v := reader.Read()
		if v == value.EOF {
			return nil
		}

		if err, ok := v.(value.Error); ok {
			return err
		}

		// Evaluate the expression.
		result := Eval(v)

		// Print the result.
		result.Write(w)
		fmt.Println()
	}
}

// Load evaluates an entire file as if in the REPL.
func Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return Run(file, os.Stdout)
}
