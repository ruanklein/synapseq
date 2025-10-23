/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
	t "github.com/ruanklein/synapseq/internal/types"
)

// generate handles the generation of a WAV file from a sequence
func generate(app *t.AppContext) error {
	if app == nil {
		return errors.New("generate: app context is nil")
	}

	sequence := app.Sequence
	if sequence == nil {
		return errors.New("generate: sequence is nil")
	}

	options := sequence.Options
	if options == nil {
		return errors.New("generate: sequence options are nil")
	}

	if !app.Quiet {
		for _, c := range sequence.Comments {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	// Create audio renderer
	renderer, err := audio.NewAudioRenderer(sequence.Periods, &audio.AudioRendererOptions{
		SampleRate:     options.SampleRate,
		Volume:         options.Volume,
		GainLevel:      options.GainLevel,
		BackgroundPath: options.BackgroundPath,
		Quiet:          app.Quiet,
		Debug:          app.Debug,
	})
	if err != nil {
		return fmt.Errorf("generate: %v", err)
	}

	if app.OutputFile == "-" {
		if app.Debug {
			return errors.New("generate: cannot use debug mode with raw output to stdout")
		}
		renderer.RenderRaw(os.Stdout)
		return nil
	}

	// Render to WAV file
	if err := renderer.RenderWav(app.OutputFile); err != nil {
		return fmt.Errorf("generate: %v", err)
	}

	// Embed text sequence metadata
	if app.Format == t.FormatText && !app.NoEmbedMetadata {
		if err := audio.WriteICMTChunkFromTextFile(app.OutputFile, app.InputFile); err != nil {
			return fmt.Errorf("generate: %v", err)
		}
	}

	return nil
}
