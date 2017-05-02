package main

import (
	"flag"
	"os"
	"testing"
)

var systemTestFlag = flag.Bool("test.system", false, "Run main as a test")

// TestMain runs a test using the main() interpreter entrypoint.
func TestMain(t *testing.T) {
	if *systemTestFlag {
		os.Stdin.Close() // interpreter exits on EOF
		main()
	}
}
