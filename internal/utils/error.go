package utils

import (
	"fmt"
	"os"
)

func Error(message string) {
	fmt.Fprintf(os.Stderr, "synapseq: %s\n", message)
	os.Exit(1)
}
