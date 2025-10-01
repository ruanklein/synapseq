/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package cli

import (
	"flag"
	"fmt"
	"os"
)

// CLIOptions holds command-line options
type CLIOptions struct {
	// Show version information and exit
	ShowVersion bool
	// Quiet mode, suppress non-error output
	Quiet bool
	// Debug mode, no wav output
	Debug bool
	// Show help message and exit
	ShowHelp bool
	// Read input as JSON format
	FormatJSON bool
}

// Usage prints the usage information
func Usage() {
	fmt.Fprintf(os.Stderr, "SynapSeq - Synapse-Sequenced Brainwave Generator, version %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "(c) 2025 Ruan, https://ruan.sh\n")
	fmt.Fprintf(os.Stderr, "Released under the GNU GPL v2. See file COPYING for details.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: synapseq [options] <input-file.spsq> <output-file.wav>\n")
}

// Help prints the help message
func Help() {
	Usage()
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -version       Show version information\n")
	fmt.Fprintf(os.Stderr, "  -json          Read input as JSON format\n")
	fmt.Fprintf(os.Stderr, "  -quiet         Enable quiet mode (suppress non-error output)\n")
	fmt.Fprintf(os.Stderr, "  -debug         Enable debug mode (no wav output)\n")
	fmt.Fprintf(os.Stderr, "  -help          Show this help message\n")
}

// ShowVersion prints the version information
func ShowVersion() {
	fmt.Printf("SynapSeq version %s\n", VERSION)
}

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.Usage = Help

	fs.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	fs.BoolVar(&opts.FormatJSON, "json", false, "Read input as JSON format")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.Debug, "debug", false, "Enable debug mode")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")

	err := fs.Parse(os.Args[1:])
	return opts, fs.Args(), err
}
