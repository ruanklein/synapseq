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

	presets, options, err := sequence.LoadSequence(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sequence Options: %+v\n\n", options)

	// Debug presets
	for _, p := range presets {
		fmt.Printf("Preset: %s\n", p.Name)
		for i, voice := range p.Voice {
			if voice.Type != t.VoiceOff {
				fmt.Printf("  Voice (%d): %+v\n", i, voice)
			}
		}
	}

}
