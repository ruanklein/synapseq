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

	"github.com/ruanklein/synapseq/internal/cli"
	"github.com/ruanklein/synapseq/internal/sequence"
	t "github.com/ruanklein/synapseq/internal/types"
)

// main is the entry point of the SynapSeq application
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

	appContext := &t.AppContext{
		Mode:            t.ModeGenerate,
		InputFile:       args[0],
		OutputFile:      args[1],
		Format:          t.FormatText,
		Quiet:           opts.Quiet,
		Debug:           opts.Debug,
		NoEmbedMetadata: opts.NoEmbedMetadata,
	}

	if opts.ExtractTextSequence {
		appContext.Mode = t.ModeExtract
	}
	if opts.ConvertToText {
		appContext.Mode = t.ModeConvert
	}

	if appContext.Mode == t.ModeExtract {
		if err = extract(appContext); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if opts.FormatJSON {
		appContext.Format = t.FormatJSON
	}
	if opts.FormatXML {
		appContext.Format = t.FormatXML
	}
	if opts.FormatYAML {
		appContext.Format = t.FormatYAML
	}

	if appContext.Format != t.FormatText {
		// Load structured sequence
		appContext.Sequence, err = sequence.LoadStructuredSequence(appContext.InputFile, appContext.Format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Load text sequence
		appContext.Sequence, err = sequence.LoadTextSequence(appContext.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
	}

	if appContext.Mode == t.ModeConvert {
		if err = convert(appContext); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err = generate(appContext); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

}
