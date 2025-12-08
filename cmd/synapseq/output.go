package main

import (
	"fmt"
	"os"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// OutputOptions defines options for processing sequence output
type outputOptions struct {
	OutputFile       string
	Quiet            bool
	Play             bool
	Mp3              bool
	UnsafeNoMetadata bool
	FFplayPath       string
	FFmpegPath       string
}

// processSequenceOutput processes the output of a loaded sequence
func processSequenceOutput(appCtx *synapseq.AppContext, opts *outputOptions) error {
	// --- Handle Stream mode (output = "-")
	if opts.OutputFile == "-" {
		return appCtx.Stream(os.Stdout)
	}

	// --- Unsafe mode
	if opts.UnsafeNoMetadata {
		var err error
		appCtx, err = appCtx.WithUnsafeNoMetadata()
		if err != nil {
			return err
		}
	}

	// --- Print comments
	if !opts.Quiet {
		for _, c := range appCtx.Comments() {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	// --- Handle Play using external ffplay
	if opts.Play {
		return externalPlay(opts.FFplayPath, appCtx)
	}

	// --- Handle MP3 output using external ffmpeg
	if opts.Mp3 {
		return externalMp3(opts.FFmpegPath, appCtx)
	}

	// Default: Render to WAV
	return appCtx.WAV()
}
