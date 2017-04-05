package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/danielwhite/microlisp/scan"
)

func main() {
	log.SetFlags(0)

	s := scan.New(os.Stdin)
	s.Error = func(s *scan.Scanner, msg string) { log.Fatal(msg) }
	for {
		// Read the next expression from the input.
		node, err := Read(s)
		if err == io.EOF {
			log.Print("end of input stream reached")
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Evaluate the expression.
		result, err := Eval(node)
		if err != nil {
			log.Print(err)
		}

		// Print the result.
		Fprint(os.Stdout, result)
		fmt.Println()
	}
}
