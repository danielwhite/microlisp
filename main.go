package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	s := NewScanner(os.Stdin)
	s.Error = func(s *Scanner, msg string) { log.Fatal(msg) }
	for {
		node, err := Read(s)
		if err == io.EOF {
			log.Print("end of input stream reached")
			break
		} else if err != nil {
			log.Fatal(err)
		}
		Fprint(os.Stdout, node)
		fmt.Println()
	}
}
