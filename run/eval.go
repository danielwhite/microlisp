// Package run provides a Lisp runtime.
package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

// UserEnvironment is an environment which inherits from the system
// environment. Definitions introduced by a user will be bound here.
//
// This exists here to provide isolation so that new definitions do
// not change the behaviour of the system environment.
var UserEnvironment value.Environment

func init() {
	Reset()
}

// Reset the environment for the runtime to an empty state.
func Reset() {
	UserEnvironment = value.NewEnv(value.SystemEnvironment)
	UserEnvironment.Define("user-environment", UserEnvironment)
}

// Eval applies rules to an expression, and returns an expression that
// is the value.
func Eval(expr value.Value) (v value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v = r.(value.Error)
		}
	}()
	v = expr.Eval(UserEnvironment)
	return
}

// EvalString evaluates the first Lisp expression in a string.
func EvalString(expr string) value.Value {
	scanner := scan.New(strings.NewReader(expr))
	reader := read.New(scanner)

	// Read the next expression from the input.
	v := reader.Read()
	if v == value.EOF {
		return nil
	}
	if err, ok := v.(value.Error); ok {
		return err
	}

	// Evaluate the expression.
	return Eval(v)
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
		fmt.Fprintln(w, result.String())
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
