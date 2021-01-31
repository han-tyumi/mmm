package utils

import (
	"fmt"
	"os"
)

// Exit causes the program to exit.
func Exit() {
	os.Exit(1)
}

// Error causes the program to exit after printing to stderr.
func Error(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	Exit()
}
