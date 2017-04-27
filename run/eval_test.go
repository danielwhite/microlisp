package run

import (
	"testing"

	"whitehouse.id.au/microlisp/value"
)

func TestEnv(t *testing.T) {
	defer Reset() // clean up environment post-test

	// Auto-quoted name of symbol; will be the result if not bound.
	sym := value.Intern("foo")

	// FIXME: Auto-quoting should be removed to simplify this
	// test. We should expect NIL (or maybe an error) for unbound
	// variables.
	if v := EvalString("foo"); v != sym {
		t.Fatalf("foo was already bound: %s", v)
	}

	// Bind a function to the symbol; and ensure it is a function.
	v := EvalString(EvalString("(defun foo () (quote x))").String())
	if _, ok := v.(value.Function); !ok {
		t.Fatalf("expected foo to be bound to function, got: %s", v)
	}

	// Clear environment, and ensure it is no longer bound.
	Reset()
	if v := EvalString("foo"); v != sym {
		t.Fatalf("foo was already bound: %s", v)
	}
}
