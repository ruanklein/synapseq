/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core_test

import (
	"fmt"
	"log"
	"os"

	synapseq "github.com/ruanklein/synapseq/core"
)

func ExampleNewAppContext() {
	// Create a new application context for text format
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav", "text")
	// To use other formats, simply change the file name and format string:
	// ctx, err := synapseq.NewAppContext("input.json", "output.wav", "json")
	// ctx, err := synapseq.NewAppContext("input.xml", "output.wav", "xml")
	// ctx, err := synapseq.NewAppContext("input.yaml", "output.wav", "yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AppContext created with format: %s\n", ctx.Format())
	// Output: AppContext created with format: text
}

func ExampleAppContext_LoadSequence() {
	// Create a new application context for text format
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("Sequence loaded successfully with format: %s\n", ctx.Format())
	// Output: Sequence loaded successfully with format: text
}

func ExampleAppContext_WAV() {
	// Create a new application context for text format
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Optional: Enable verbose output
	// Replace with an io.Writer, e.g., os.Stderr
	ctx = ctx.WithVerbose(os.Stderr)

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Generate the WAV file
	// if err := ctx.WAV(); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("WAV file generated successfully with format: %s\n", ctx.Format())
	// Output: WAV file generated successfully with format: text
}

func ExampleAppContext_Stream() {
	// Create a new application context for text format
	ctx, err := synapseq.NewAppContext("input.spsq", "", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Stream the RAW data to standard output (44100 Hz [default], 24-bit, stereo)
	// Replace with an io.Writer, e.g., os.Stdout
	// if err := ctx.Stream(os.Stdout); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("RAW data streamed successfully with format: %s\n", ctx.Format())
	// Output: RAW data streamed successfully with format: text
}

func ExampleAppContext_Comments() {
	// Create a new application context for text format
	ctx, err := synapseq.NewAppContext("input.spsq", "output.wav", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Retrieve comments from the sequence
	// for _, comment := range ctx.Comments() {
	//	fmt.Println(comment)
	// }

	fmt.Printf("Comments retrieved successfully with format: %s\n", ctx.Format())
	// Output: Comments retrieved successfully with format: text
}

func ExampleAppContext_Text() {
	// Create a new application context for JSON format
	ctx, err := synapseq.NewAppContext("input.json", "", "json")
	if err != nil {
		log.Fatal(err)
	}

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Convert the sequence to text format
	// content, err := ctx.Text()
	// if err != nil {
	//	log.Fatal(err)
	// }
	// fmt.Println(content)

	fmt.Printf("Sequence converted to text format successfully from format: %s\n", ctx.Format())
	// Output: Sequence converted to text format successfully from format: json
}

func ExampleAppContext_SaveText() {
	// Create a new application context for XML format
	ctx, err := synapseq.NewAppContext("input.xml", "output.spsq", "xml")
	if err != nil {
		log.Fatal(err)
	}

	// Load the sequence
	// if err := ctx.LoadSequence(); err != nil {
	//	log.Fatal(err)
	// }

	// Save the sequence as text format
	// if err := ctx.SaveText(); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Printf("Sequence saved as text format successfully from format: %s\n", ctx.Format())
	// Output: Sequence saved as text format successfully from format: xml
}

func ExampleExtract() {
	// Extract text sequence from WAV file
	// content, err := synapseq.Extract("input.wav")
	// if err != nil {
	//	log.Fatal(err)
	// }
	// fmt.Println(content)

	fmt.Println("Text sequence extracted successfully from WAV file.")
	// Output: Text sequence extracted successfully from WAV file.
}

func ExampleSaveExtracted() {
	// Save extracted text sequence from WAV file to output file
	// if err := synapseq.SaveExtracted("input.wav", "output.spsq"); err != nil {
	//	log.Fatal(err)
	// }

	fmt.Println("Text sequence extracted and saved successfully from WAV file.")
	// Output: Text sequence extracted and saved successfully from WAV file.
}
