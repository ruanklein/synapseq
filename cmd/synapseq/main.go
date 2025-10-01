/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
	"github.com/ruanklein/synapseq/internal/cli"
	"github.com/ruanklein/synapseq/internal/sequence"
)

func main() {
	opts, args, err := cli.ParseFlags()
	if err != nil {
		os.Exit(2)
	}

	if opts.ShowHelp {
		cli.Help()
		os.Exit(0)
	}
	if opts.ShowVersion {
		cli.ShowVersion()
		os.Exit(0)
	}

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "synapseq: invalid number of arguments\n")
		cli.Usage()
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	var result *sequence.LoadResult
	if opts.FormatJSON {
		// Load JSON sequence
		var err error
		result, err = sequence.LoadJSONSequence(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load text sequence
		var err error
		result, err = sequence.LoadTextSequence(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
	}

	if !opts.Quiet {
		for _, c := range result.Comments {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	options := result.Options

	// Create audio renderer
	renderer, err := audio.NewAudioRenderer(result.Periods, &audio.AudioRendererOptions{
		SampleRate:     options.SampleRate,
		Volume:         options.Volume,
		GainLevel:      options.GainLevel,
		BackgroundPath: options.BackgroundPath,
		Quiet:          opts.Quiet,
		Debug:          opts.Debug,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	if outputFile == "-" {
		if opts.Debug {
			fmt.Fprintf(os.Stderr, "synapseq: cannot use debug mode with raw output to stdout\n")
			os.Exit(1)
		}
		renderer.RenderRaw(os.Stdout)
		return
	}

	// Render to WAV file
	if err := renderer.RenderWav(outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}
}
