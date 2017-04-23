package run

import (
	"strings"

	"whitehouse.id.au/microlisp/read"
	"whitehouse.id.au/microlisp/scan"
	"whitehouse.id.au/microlisp/value"
)

// ReadString reads the first textual Lisp expression from the text in
// string.
func ReadString(text string) value.Value {
	scanner := scan.New(strings.NewReader(text))
	reader := read.New(scanner)
	return reader.Read()
}
