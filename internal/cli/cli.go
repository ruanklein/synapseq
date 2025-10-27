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
	"runtime"

	"github.com/ruanklein/synapseq/v3/internal/info"
)

// CLIOptions holds command-line options
type CLIOptions struct {
	// Show version information and exit
	ShowVersion bool
	// Quiet mode, suppress non-error output
	Quiet bool
	// Test mode, validate syntax without generating output
	Test bool
	// Show help message and exit
	ShowHelp bool
	// Read input as JSON format
	FormatJSON bool
	// Read input as XML format
	FormatXML bool
	// Read input as YAML format
	FormatYAML bool
	// Extract text sequence from WAV file
	ExtractTextSequence bool
	// Do not embed metadata in output WAV file
	UnsafeNoMetadata bool
	// Convert to text from json/xml/yaml
	ConvertToText bool
}

// Help prints the help message
func Help() {
	fmt.Printf("SynapSeq - Synapse-Sequenced Brainwave Generator, version %s\n", info.VERSION)
	fmt.Printf("(c) 2025 %s, %s\n", info.AUTHOR, info.AUTHOR_URL)
	fmt.Printf("Released under the GNU GPL v2. See file COPYING for details.\n\n")

	fmt.Printf("Usage: synapseq [options] <input> <output>\n\n")

	fmt.Printf("INPUT formats:\n")
	fmt.Printf("    Local file path:     path/to/sequence.spsq\n")
	fmt.Printf("    Standard input:      -\n")
	fmt.Printf("    HTTP/HTTPS URL:      https://example.com/sequence.spsq\n")
	fmt.Printf("    Structured formats:  Use -json, -xml, or -yaml flags\n\n")

	fmt.Printf("OUTPUT formats:\n")
	fmt.Printf("    WAV file:            path/to/output.wav\n")
	fmt.Printf("    Standard output:     - (raw PCM, 24-bit stereo)\n\n")

	fmt.Printf("Options:\n")
	fmt.Printf("  -json          		Read input as JSON format\n")
	fmt.Printf("  -xml           		Read input as XML format\n")
	fmt.Printf("  -yaml          		Read input as YAML format\n")
	fmt.Printf("  -quiet         		Suppress non-error output\n")
	fmt.Printf("  -test          		Validate syntax without generating output\n")
	fmt.Printf("  -extract       		Extract text sequence from WAV file\n")
	fmt.Printf("  -convert       		Convert to text from json/xml/yaml\n")
	fmt.Printf("  -unsafe-no-metadata  	  	Do not embed metadata in output WAV file\n")
	fmt.Printf("  -version       		Show version information\n")
	fmt.Printf("  -help         		Show this help message\n\n")

	fmt.Printf("Examples:\n")
	fmt.Printf("  synapseq sequence.spsq output.wav\n")
	fmt.Printf("  synapseq -test sequence.spsq output.wav\n")
	fmt.Printf("  synapseq -json sequence.json output.wav\n")
	fmt.Printf("  cat sequence.spsq | synapseq - output.wav\n")
	fmt.Printf("  synapseq https://example.com/sequence.spsq output.wav\n")
	fmt.Printf("  synapseq sequence.spsq - | play -t raw -r 44100 -e signed-integer -b 24 -c 2 -\n")
	fmt.Printf("  synapseq -extract sequence.wav output.spsq\n")
	fmt.Printf("  synapseq -convert -json sequence.json output.spsq\n\n")

	fmt.Printf("For detailed documentation:\n")
	if runtime.GOOS != "windows" {
		fmt.Printf("  man synapseq\n")
	}
	fmt.Printf("  %s\n", info.REPOSITORY)
}

// ShowVersion prints the version information
func ShowVersion() {
	fmt.Printf("SynapSeq %s (%s) built %s for %s/%s\n",
		info.VERSION,
		info.GIT_COMMIT,
		info.BUILD_DATE,
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// ParseFlags parses command-line flags and returns CLIOptions
func ParseFlags() (*CLIOptions, []string, error) {
	opts := &CLIOptions{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use -help flag for usage information.\n")
	}

	fs.BoolVar(&opts.ShowVersion, "version", false, "Show version information")
	fs.BoolVar(&opts.FormatJSON, "json", false, "Read input as JSON format")
	fs.BoolVar(&opts.FormatXML, "xml", false, "Read input as XML format")
	fs.BoolVar(&opts.FormatYAML, "yaml", false, "Read input as YAML format")
	fs.BoolVar(&opts.Quiet, "quiet", false, "Enable quiet mode")
	fs.BoolVar(&opts.Test, "test", false, "Validate syntax without generating output")
	fs.BoolVar(&opts.ExtractTextSequence, "extract", false, "Extract text sequence from WAV file")
	fs.BoolVar(&opts.UnsafeNoMetadata, "unsafe-no-metadata", false, "Do not embed metadata in output WAV file")
	fs.BoolVar(&opts.ConvertToText, "convert", false, "Convert to text from json/xml/yaml")
	fs.BoolVar(&opts.ShowHelp, "help", false, "Show help")

	err := fs.Parse(os.Args[1:])
	return opts, fs.Args(), err
}
