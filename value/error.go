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

// Equal implements the Value interface.
func (e Error) Equal(v Value) Value {
	e2, ok := v.(Error)
	if !ok {
		return NIL
	}
	if e != e2 {
		return NIL
	}
	return T
}
