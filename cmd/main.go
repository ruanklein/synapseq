package main

import (
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
	"github.com/ruanklein/synapseq/internal/sequence"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.spsq> <output.wav>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	periods, options, err := sequence.LoadSequence(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Sequence Options:\n  %+v\n\n", options.String())

	// Debug periods
	// for _, p := range periods {
	// 	fmt.Printf("- %s\n", p.TimeString())
	// 	for _, voice := range p.VoiceStart {
	// 		if voice.Type != t.VoiceOff {
	// 			fmt.Printf("\t%s\n", voice.String())
	// 		}
	// 	}
	// }

	// Create audio renderer
	renderer, err := audio.NewAudioRenderer(periods, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating renderer: %v\n", err)
		os.Exit(1)
	}

	// Render to WAV file
	if err := renderer.RenderToWAV(outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering audio: %v\n", err)
		os.Exit(1)
	}
}
