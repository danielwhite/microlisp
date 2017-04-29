package value

import (
	"fmt"
)

// Error is a value used to represent runtime errors.
type Error string

// Errorf raises an error with a formatted message.
func Errorf(format string, a ...interface{}) {
	panic(Error(fmt.Sprintf(format, a...)))
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
