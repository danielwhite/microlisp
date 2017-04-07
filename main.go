package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/run"
	"github.com/danielwhite/microlisp/scan"
	"github.com/danielwhite/microlisp/value"
)

func main() {
	log.SetFlags(0)

	scanner := scan.New(bufio.NewReader(os.Stdin))
	reader := read.New(scanner)
	for {
		// Read the next expression from the input.
		v := reader.Read()
		if v == value.EOF {
			log.Print("end of input stream reached")
			break
		}

		if err, ok := v.(value.Error); ok {
			log.Fatal(err)
		}

		// Evaluate the expression.
		result := run.Eval(v)

		// Print the result.
		result.Write(os.Stdout)
		fmt.Println()
	}
}
