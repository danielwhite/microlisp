package value

import (
	"fmt"
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

// Errorf returns an error with a formatted message.
func Errorf(format string, a ...interface{}) Error {
	return Error(fmt.Sprintf(format, a...))
}

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}

func (e Error) String() string {
	return fmt.Sprintf("#[error: %s]", string(e))
}

// Eval implements the Value interface.
func (e Error) Eval(env Environment) Value {
	return e
}

// Equal implments the Value interface.
func (e Error) Equal(Value) Value {
	return NIL
}
