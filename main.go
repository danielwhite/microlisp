package main

import (
	"flag"
	"log"
	"os"

	"whitehouse.id.au/microlisp/run"
)

var (
	loadFlag = flag.String("load", "", "Load a file as if in the REPL")

	// TODO(danielwhite): Teach -q flag to suppress results of evaluation.
	quietFlag = flag.Bool("q", false, "Suppress the REPL prompt")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	// If a file is specified, load before proceeding.
	if *loadFlag != "" {
		if err := run.Load(*loadFlag); err != nil {
			// Typical interpreter behaviour is to allow
			// for recovery in the event of a load error.
			log.Printf("Unable to load file %q: %s", *loadFlag, err)
		}
	}

	if *quietFlag {
		run.DefaultPrompt = ""
	}

	// Enter the interactive REPL.
	run.Run(os.Stdin, os.Stdout)
}
