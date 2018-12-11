package run

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

// DefaultPrompt is printed by [Run] to prompt for input.
var DefaultPrompt = "> "

// Run a REPL loop, reading expressions from a reader, and writing the
// evaluated values to a writer.
func Run(r io.Reader, w io.Writer) error {
	return run(r, w, DefaultPrompt)
}

func run(r io.Reader, w io.Writer, prompt string) error {
	scanner := scan.New(bufio.NewReader(r))
	reader := read.New(scanner)
	for {
		io.WriteString(w, prompt)

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

	return run(file, os.Stdout, "")
}
