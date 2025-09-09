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
	// Show help message and exit
	ShowHelp bool
}

// Usage prints the usage information
func Usage() {
	fmt.Fprintf(os.Stderr, "SynapSeq - Synapse-Sequenced Brainwave Generator, version %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "(c) 2025 Ruan, https://ruan.sh\n")
	fmt.Fprintf(os.Stderr, "Released under the GNU GPL v2. See file COPYING for details.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s <input.spsq> <output.wav>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -version	Show version information\n")
	fmt.Fprintf(os.Stderr, "  -quiet    	Enable debug mode for sequences\n")
	fmt.Fprintf(os.Stderr, "  -help 	Show help\n")
}

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() *CLIOptions {
	opts := &CLIOptions{}
	flag.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	flag.BoolVar(&opts.ShowHelp, "help", false, "Show help")

	flag.Parse()
	return opts
}
