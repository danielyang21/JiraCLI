package ui

import (
	"fmt"
	"os"
)

func FatalError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func FatalIfError(err error, message string) {
	if err != nil {
		FatalError("%s: %v", message, err)
	}
}
