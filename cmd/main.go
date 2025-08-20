package main

import (
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/sequence"
	"github.com/ruanklein/synapseq/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: synapseq <file>")
		os.Exit(1)
	}

	// Debug sequence
	if err := sequence.LoadSequence(os.Args[1]); err != nil {
		utils.Error(err.Error())
	}

}
