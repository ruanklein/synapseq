package main

import (
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/sequence"
	t "github.com/ruanklein/synapseq/internal/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: synapseq <file>")
		os.Exit(1)
	}

	// waveTables := audio.InitWaveformTables()

	periods, options, err := sequence.LoadSequence(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sequence Options:\n  %+v\n\n", options.String())

	// Debug periods
	for _, p := range periods {
		fmt.Printf("Period start at %s\n", p.TimeString())
		for i, voice := range p.VoiceStart {
			if voice.Type != t.VoiceOff {
				fmt.Printf("  Voice (%d): %s\n", i+1, voice.String())
			}
		}
	}

}
