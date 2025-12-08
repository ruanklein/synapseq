//go:build !wasm

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
	"path/filepath"
	"strings"

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

	// --hub-update
	if opts.HubUpdate {
		return hubRunUpdate(opts.Quiet)
	}

	// --hub-clean
	if opts.HubClean {
		return hubRunClean(opts.Quiet)
	}

	// --hub-get
	if opts.HubGet != "" {
		var outputFile string
		if len(args) == 1 {
			outputFile = args[0]
		}
		return hubRunGet(opts.HubGet, outputFile, opts)
	}

	// --hub-list
	if opts.HubList {
		return hubRunList()
	}

	// --hub-search
	if opts.HubSearch != "" {
		return hubRunSearch(opts.HubSearch)
	}

	// --hub-download
	if opts.HubDownload != "" {
		targetDir := ""
		if len(args) == 1 {
			targetDir = args[0]
		}
		return hubRunDownload(opts.HubDownload, targetDir, opts.Quiet)
	}

	// --hub-info
	if opts.HubInfo != "" {
		return hubRunInfo(opts.HubInfo)
	}

	// --install-file-association (Windows only)
	if opts.InstallFileAssociation {
		return installWindowsFileAssociation(opts.Quiet)
	}

	// --uninstall-file-association (Windows only)
	if opts.UninstallFileAssociation {
		return uninstallWindowsFileAssociation(opts.Quiet)
	}

	// --help or missing args
	if opts.ShowHelp || len(args) == 0 {
		cli.Help()
		return nil
	}

	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid number of flags\nUse -help for usage information")
	}

	// Determine output format
	outputFormat := "wav"
	if opts.Mp3 {
		outputFormat = "mp3"
	}

	inputFile := args[0]
	outputFile := getDefaultOutputFile(inputFile, outputFormat)
	if len(args) == 2 {
		outputFile = args[1]
	}

	// --- Handle Extract mode
	if opts.ExtractTextSequence {
		if opts.Mp3 {
			if outputFile == "-" {
				content, err := externalExtractTextSequence(opts.FFprobePath, inputFile)
				if err != nil {
					return fmt.Errorf("failed to extract text sequence. Error\n  %w", err)
				}
				fmt.Println(content)
				return nil
			}

			outputFile = getDefaultOutputFile(inputFile, "spsq")
			if err := externalSaveExtractedTextSequence(opts.FFprobePath, inputFile, outputFile); err != nil {
				return fmt.Errorf("failed to extract text sequence. Error\n  %w", err)
			}

			if !opts.Quiet {
				fmt.Println("Extraction completed successfully.")
			}

			return nil
		}

		if outputFile == "-" {
			content, err := synapseq.Extract(inputFile)
			if err != nil {
				return fmt.Errorf("failed to extract text sequence. Error\n  %w", err)
			}
			fmt.Println(content)
			return nil
		}

		outputFile = getDefaultOutputFile(inputFile, "spsq")
		if err := synapseq.SaveExtracted(inputFile, outputFile); err != nil {
			return fmt.Errorf("failed to extract text sequence. Error\n  %w", err)
		}

		if !opts.Quiet {
			fmt.Println("Extraction completed successfully.")
		}
		return nil
	}

	// Detect format flags
	format := detectFormat(opts)

	appCtx, err := synapseq.NewAppContext(inputFile, outputFile, format)
	if err != nil {
		return err
	}

	if !opts.Quiet && outputFile != "-" {
		appCtx = appCtx.WithVerbose(os.Stderr)
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
				return fmt.Errorf("failed to convert to text. Error\n  %w", err)
			}
			fmt.Println(content)
			return nil
		}

		if err := appCtx.SaveText(); err != nil {
			return fmt.Errorf("failed to convert to text. Error\n  %w", err)
		}

		if !opts.Quiet {
			fmt.Println("Conversion completed successfully.")
		}
		return nil
	}

	// --- Process output using centralized handler
	outputOpts := &outputOptions{
		OutputFile:       outputFile,
		Quiet:            opts.Quiet,
		Play:             opts.Play,
		Mp3:              opts.Mp3,
		UnsafeNoMetadata: opts.UnsafeNoMetadata,
		FFplayPath:       opts.FFplayPath,
		FFmpegPath:       opts.FFmpegPath,
	}

	return processSequenceOutput(appCtx, outputOpts)
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

// getDefaultOutputFile generates a default output filename based on the input filename
func getDefaultOutputFile(inputFile string, extension string) string {
	base := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	return base + "." + extension
}
