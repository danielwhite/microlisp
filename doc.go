/*
Microlisp is a Lisp interpreter. It's primarily a toy, and is a work in progress.

Its initial definition is based on McCarthy's "A Micro-Manual for Lisp: Not the Whole Truth".

# Environments

An environment is a value that maintains mappings between a symbol and
a value. If an environment has a parent, then bindings are inherited.

The evaluator typically provides an environment implicitly.

Variables.

	system-environment	Primitives are bound in this environment.
	user-environment	User definitions are bound here, and it inherits system-environment.

Functions.

	environment-bindings	The environment's bindings represented as an association list.

# Errors

The error system is very simple. If an error value is thrown, it stops
the current evaluation, and is handled and printed by the REPL.

# Functions

	error		Raise an error value with a message composed of its arguments.
	ignore-errors	Invoke a function, trapping any errors thrown as a return value.
*/
package main // import "whitehouse.id.au/microlisp"
