/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	"fmt"
	"io"

	"github.com/ruanklein/synapseq/internal/audio"
	"github.com/ruanklein/synapseq/internal/info"
	t "github.com/ruanklein/synapseq/internal/types"
)

// generate generates the audio renderer based on the loaded sequence
func (ac *AppContext) generate(debug bool) (*audio.AudioRenderer, error) {
	sequence := ac.sequence
	if sequence == nil {
		return nil, fmt.Errorf("sequence is nil")
	}

	options := sequence.Options
	if options == nil {
		return nil, fmt.Errorf("sequence options are nil")
	}

	renderer, err := audio.NewAudioRenderer(sequence.Periods, &audio.AudioRendererOptions{
		SampleRate:     options.SampleRate,
		Volume:         options.Volume,
		GainLevel:      options.GainLevel,
		BackgroundPath: options.BackgroundPath,
		Quiet:          !ac.verbose,
		Debug:          debug,
	})
	if err != nil {
		return nil, err
	}

	return renderer, nil
}

// WAV generates the WAV file from the loaded sequence
func (ac *AppContext) WAV() error {
	renderer, err := ac.generate(false)
	if err != nil {
		return err
	}

	err = renderer.RenderWav(ac.outputFile)
	if err != nil {
		return err
	}

	if ac.format == t.FormatText && !ac.unsafeNoMetadata {
		metadata, err := info.NewMetadata(ac.inputFile)
		if err != nil {
			return err
		}

		if err = audio.WriteICMTChunkFromTextFile(ac.outputFile, metadata); err != nil {
			return err
		}

	}

	return nil
}

// Stream generates the raw audio stream from the loaded sequence
func (ac *AppContext) Stream(data io.Writer) error {
	renderer, err := ac.generate(false)
	if err != nil {
		return err
	}

	err = renderer.RenderRaw(data)
	if err != nil {
		return err
	}

	return nil
}

// Debug runs the audio generation in debug mode
func (ac *AppContext) Debug() error {
	_, err := ac.generate(true)
	if err != nil {
		return err
	}

	return nil
}
