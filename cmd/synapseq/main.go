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
		cli.Help()
		os.Exit(1)
	}

	inputFile := args[0]
	outputFile := args[1]

	if opts.ExtractTextSequence {
		// Extract text sequence from WAV file
		if err := audio.ExtractTextSequenceFromWAV(inputFile, outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
		if !opts.Quiet {
			fmt.Fprintf(os.Stderr, "Extracted text sequence to %s\n", outputFile)
		}
		return
	}

	// Default to text sequence
	format := "text"
	if opts.FormatJSON {
		format = "json"
	}
	if opts.FormatXML {
		format = "xml"
	}
	if opts.FormatYAML {
		format = "yaml"
	}

	var result *sequence.LoadResult
	if format != "text" {
		// Load structured sequence
		var err error
		result, err = sequence.LoadStructuredSequence(inputFile, format)
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

	// Embed text sequence metadata
	if format == "text" && !opts.NoEmbedMetadata {
		if err := audio.WriteICMTChunkFromTextFile(outputFile, inputFile, format); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			fmt.Fprintf(os.Stderr, "Try running again with --no-embed to skip embedding metadata.\n")
			os.Exit(1)
		}
	}
}
