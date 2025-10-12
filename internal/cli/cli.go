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
	// Read input as XML format
	FormatXML bool
	// Read input as YAML format
	FormatYAML bool
}

// Help prints the help message
func Help() {
	fmt.Fprintf(os.Stderr, "SynapSeq - Synapse-Sequenced Brainwave Generator, version %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "(c) 2025 Ruan, https://ruan.sh\n")
	fmt.Fprintf(os.Stderr, "Released under the GNU GPL v2. See file COPYING for details.\n\n")

	fmt.Fprintf(os.Stderr, "Usage: synapseq [options] <input> <output>\n\n")

	fmt.Fprintf(os.Stderr, "INPUT formats:\n")
	fmt.Fprintf(os.Stderr, "    Local file path:     path/to/sequence.spsq\n")
	fmt.Fprintf(os.Stderr, "    Standard input:      -\n")
	fmt.Fprintf(os.Stderr, "    HTTP/HTTPS URL:      https://example.com/sequence.spsq\n")
	fmt.Fprintf(os.Stderr, "    Structured formats:  Use -json, -xml, or -yaml flags\n\n")

	fmt.Fprintf(os.Stderr, "OUTPUT formats:\n")
	fmt.Fprintf(os.Stderr, "    WAV file:            path/to/output.wav\n")
	fmt.Fprintf(os.Stderr, "    Standard output:     - (raw PCM, 24-bit stereo)\n\n")

	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -json          Read input as JSON format\n")
	fmt.Fprintf(os.Stderr, "  -xml           Read input as XML format\n")
	fmt.Fprintf(os.Stderr, "  -yaml          Read input as YAML format\n")
	fmt.Fprintf(os.Stderr, "  -quiet         Suppress non-error output\n")
	fmt.Fprintf(os.Stderr, "  -debug         Validate syntax without generating output\n")
	fmt.Fprintf(os.Stderr, "  -version       Show version information\n")
	fmt.Fprintf(os.Stderr, "  -help          Show this help message\n\n")

	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  synapseq sequence.spsq output.wav\n")
	fmt.Fprintf(os.Stderr, "  synapseq -json sequence.json output.wav\n")
	fmt.Fprintf(os.Stderr, "  cat sequence.spsq | synapseq - output.wav\n")
	fmt.Fprintf(os.Stderr, "  synapseq https://example.com/sequence.spsq output.wav\n")
	fmt.Fprintf(os.Stderr, "  synapseq sequence.spsq - | play -t raw -r 44100 -e signed-integer -b 24 -c 2 -\n\n")

	fmt.Fprintf(os.Stderr, "For detailed documentation, visit: https://github.com/ruanklein/synapseq\n")
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
	fs.BoolVar(&opts.FormatXML, "xml", false, "Read input as XML format")
	fs.BoolVar(&opts.FormatYAML, "yaml", false, "Read input as YAML format")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.Debug, "debug", false, "Enable debug mode")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")

	err := fs.Parse(os.Args[1:])
	return opts, fs.Args(), err
}
