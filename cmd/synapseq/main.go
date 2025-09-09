package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
	"github.com/ruanklein/synapseq/internal/cli"
	"github.com/ruanklein/synapseq/internal/sequence"
)

func main() {
	opts := cli.ParseFlags()

	if opts.ShowHelp {
		cli.Usage()
		os.Exit(1)
	}
	if opts.ShowVersion {
		fmt.Printf("SynapSeq version %s\n", cli.VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "synapseq: invalid number of arguments\n")
		cli.Usage()
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	// Load sequence
	periods, options, err := sequence.LoadTextSequence(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	// Create audio renderer
	renderer, err := audio.NewAudioRenderer(periods, &audio.AudioRendererOptions{
		SampleRate:     options.SampleRate,
		Volume:         options.Volume,
		GainLevel:      options.GainLevel,
		BackgroundPath: options.BackgroundPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	// Render to WAV file
	if err := renderer.RenderToWAV(outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}
}
