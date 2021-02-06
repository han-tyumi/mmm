package utils

import (
	"fmt"
	"os"
)

// Error causes the program to exit after printing to stderr.
func Error(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
