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

	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/internal/cli"
)

// main is the entry point of the SynapSeq application
func main() {
	opts, args, err := cli.ParseFlags()
	if err != nil {
		os.Exit(1)
	}

	if err := run(opts, args); err != nil {
		fmt.Fprintf(os.Stderr, "synapseq: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main application logic based on CLI options and arguments
func run(opts *cli.CLIOptions, args []string) error {
	// --version
	if opts.ShowVersion {
		cli.ShowVersion()
		return nil
	}

	// --help or missing args
	if opts.ShowHelp || len(args) == 0 {
		cli.Help()
		return nil
	}

	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid number of flags\nUse -help for usage information")
	}

	inputFile := args[0]
	outputFile := "-"
	if len(args) == 2 {
		outputFile = args[1]
	}

	// --- Handle Extract mode
	if opts.ExtractTextSequence {
		if outputFile == "-" {
			content, err := synapseq.Extract(inputFile)
			if err != nil {
				return fmt.Errorf("failed to extract text sequence: %w", err)
			}
			fmt.Println(content)
			return nil
		}

		if err := synapseq.SaveExtracted(inputFile, outputFile); err != nil {
			return fmt.Errorf("failed to extract text sequence: %w", err)
		}

		fmt.Println("Extraction completed successfully.")
		return nil
	}

	// Detect format flags
	format := detectFormat(opts)

	appCtx, err := synapseq.NewAppContext(inputFile, outputFile, format)
	if err != nil {
		return err
	}

	if !opts.Quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stdout)
	}

	// Load sequence file
	if err := appCtx.LoadSequence(); err != nil {
		return err
	}

	// --- Handle Test mode (no output required)
	if opts.Test {
		if !opts.Quiet {
			fmt.Println("Sequence is valid.")
		}
		return nil
	}

	// --- Handle Convert mode
	if opts.ConvertToText {
		if outputFile == "-" {
			content, err := appCtx.Text()
			if err != nil {
				return fmt.Errorf("failed to convert to text: %w", err)
			}
			fmt.Println(content)
			return nil
		}

		if err := appCtx.SaveText(); err != nil {
			return fmt.Errorf("failed to convert to text: %w", err)
		}

		fmt.Println("Conversion completed successfully.")
		return nil
	}

	// --- Handle Stream mode (output = "-")
	if outputFile == "-" {
		return appCtx.Stream(os.Stdout)
	}

	// --- Unsafe mode
	if opts.UnsafeNoMetadata {
		appCtx, err = appCtx.WithUnsafeNoMetadata()
		if err != nil {
			return err
		}
	}

	// --- Render to WAV (default mode)
	if !opts.Quiet {
		for _, c := range appCtx.Comments() {
			fmt.Printf("> %s\n", c)
		}
	}

	return appCtx.WAV()
}

// detectFormat detects the input format based on CLI options
func detectFormat(opts *cli.CLIOptions) string {
	switch {
	case opts.FormatJSON:
		return "json"
	case opts.FormatXML:
		return "xml"
	case opts.FormatYAML:
		return "yaml"
	default:
		return "text"
	}
}
