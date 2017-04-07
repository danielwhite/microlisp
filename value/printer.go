package value

import (
	"bytes"
)

// Sprint prints a Lisp value to a string.
func Sprint(v Value) string {
	var buf bytes.Buffer
	v.Write(&buf)
	return buf.String()
}

type printer struct {
	buf bytes.Buffer
}
