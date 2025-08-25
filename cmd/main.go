package main

import (
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
	"github.com/ruanklein/synapseq/internal/sequence"
	"github.com/ruanklein/synapseq/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: synapseq <file>")
		os.Exit(1)
	}

	// waveTables := audio.InitWaveformTables()

	presets, options, err := sequence.LoadSequence(os.Args[1])
	if err != nil {
		utils.Error(err.Error())
	}

	fmt.Printf("Sequence Options: %+v\n\n", options)

	// Debug presets
	for _, p := range presets {
		fmt.Printf("Preset: %s\n", p.Name)
		for i, voice := range p.Voice {
			if voice.Type != audio.VoiceOff {
				fmt.Printf("  Voice (%d): %+v\n", i, voice)
			}
		}
	}

}
