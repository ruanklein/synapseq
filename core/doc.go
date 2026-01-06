/*
Package core provides the application context and core functionality
for the SynapSeq brainwave generator.

# Overview

This package is designed to be used as a library by other Go projects
that want to integrate SynapSeq audio generation capabilities.

# Supported Formats

SynapSeq supports multiple input formats:
  - text (.spsq): Human-readable text format with presets
  - json: Structured JSON format
  - xml: Structured XML format
  - yaml: Structured YAML format

# Example Usage

	package main

	import (
	    "log"
	    "os"

	    synapseq "github.com/synapseq-foundation/synapseq/v3/core"
	)

	func main() {
	    // Create application context
	    ctx, err := synapseq.NewAppContext("input.spsq", "output.wav", "text")
	    if err != nil {
	        log.Fatal(err)
	    }

	    // Enable verbose output (optional)
	    ctx = ctx.WithVerbose(os.Stderr)

		// Load sequence (required before generating WAV, streaming, or converting)
		if err := ctx.LoadSequence(); err != nil {
			log.Fatal(err)
		}

	    // Generate WAV file
	    if err := ctx.WAV(); err != nil {
	        log.Fatal(err)
	    }
	}

# File Paths

Input and output files support:
  - Local file paths: "path/to/file.spsq"
  - Standard input: "-" (only for input files)
  - HTTP/HTTPS URLs: "https://example.com/sequence.spsq"

# Thread Safety

AppContext methods are safe for concurrent use as they return new instances
rather than modifying the original context.

# More Information

For complete documentation and examples, see:
https://github.com/synapseq-foundation/synapseq
*/
package core
