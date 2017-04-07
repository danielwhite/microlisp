package run

import (
	"strings"

	"github.com/danielwhite/microlisp/read"
	"github.com/danielwhite/microlisp/scan"
	"github.com/danielwhite/microlisp/value"
)

// ReadString reads the first textual Lisp expression from the text in
// string.
func ReadString(text string) value.Value {
	scanner := scan.New(strings.NewReader(text))
	reader := read.New(scanner)
	return reader.Read()
}
