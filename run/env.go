package run

import "whitehouse.id.au/microlisp/value"

// UserEnvironment is an environment which inherits from the system
// environment. Definitions introduced by a user will be bound here.
//
// This exists here to provide isolation so that new definitions do
// not change the behaviour of the system environment.
var UserEnvironment value.Environment

func init() {
	Reset()
}

// Reset the environment for the runtime to an empty state.
func Reset() {
	UserEnvironment = value.NewEnv(value.SystemEnvironment)
	UserEnvironment.Define("user-environment", UserEnvironment)
}
