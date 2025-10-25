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

	"github.com/ruanklein/synapseq/core"
	"github.com/ruanklein/synapseq/internal/cli"
)

// main is the entry point of the SynapSeq application
func main() {
	opts, args, err := cli.ParseFlags()
	if err != nil {
		os.Exit(1)
	}

	if opts.ShowHelp {
		cli.Help()
		return
	}
	if opts.ShowVersion {
		cli.ShowVersion()
		return
	}

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "synapseq: invalid number of arguments\n")
		cli.Help()
		os.Exit(1)
	}

	var (
		inputFile  = args[0]
		outputFile = args[1]
		format     = "text"
	)

	if opts.FormatJSON {
		format = "json"
	}
	if opts.FormatXML {
		format = "xml"
	}
	if opts.FormatYAML {
		format = "yaml"
	}

	appContext, err := core.NewAppContext(inputFile, outputFile, format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	if !opts.Quiet && outputFile != "-" {
		appContext = appContext.WithVerbose(os.Stderr)
	}

	if opts.ExtractTextSequence {
		if outputFile == "-" {
			content, err := appContext.Extract()
			if err != nil {
				fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(content)
			return
		}

		if err = appContext.SaveExtracted(); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}

		if !opts.Quiet {
			fmt.Println("Extraction completed successfully.")
		}
		return
	}

	if err = appContext.LoadSequence(); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}

	if opts.Test {
		if !opts.Quiet {
			fmt.Println("Sequence is valid.")
		}
		return
	}

	if opts.ConvertToText {
		if outputFile == "-" {
			content, err := appContext.Text()
			if err != nil {
				fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(content)
			return
		}

		if err = appContext.SaveText(); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}

		if !opts.Quiet {
			fmt.Println("Conversion completed successfully.")
		}
		return
	}

	if outputFile == "-" {
		if err = appContext.Stream(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if !opts.Quiet {
		for _, comment := range appContext.Comments() {
			fmt.Printf("> %s\n", comment)
		}
	}

	if opts.UnsafeNoMetadata {
		appContext, err = appContext.WithUnsafeNoMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
			os.Exit(1)
		}
	}

	if err = appContext.WAV(); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}
}
