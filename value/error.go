package value

import (
	"fmt"
	"io"
)

// Panic raises an error with the given message.
func Panic(msg string) {
	panic(Error(msg))
}

// Panicf raises an error with a formatted message.
func Panicf(format string, a ...interface{}) {
	Panic(fmt.Sprintf(format, a...))
}

// Error is a value used to represent runtime errors.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

// Write implements the Value interface.
func (e Error) Write(w io.Writer) {
	fmt.Fprintf(w, "#[error: %s]", e)
}

// Eval implements the Value interface.
func (e Error) Eval(env Environment) Value {
	panic(e) // why are we evaluating an error!?
}

// Equal implments the Value interface.
func (e Error) Equal(Value) Value {
	return NIL
}
