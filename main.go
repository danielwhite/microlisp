package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	s := NewScanner(os.Stdin)
	s.Error = func(s *Scanner, msg string) { log.Fatal(msg) }
	for {
		tok, lit := s.Scan()
		if tok == EOF {
			break
		}
		fmt.Println(tok, lit)
	}
}
