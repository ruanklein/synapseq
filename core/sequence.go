//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	seq "github.com/ruanklein/synapseq/v3/internal/sequence"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// LoadSequence loads the sequence from the input file based on the specified format
func (ac *AppContext) LoadSequence() error {
	var err error
	if ac.format == t.FormatText {
		ac.sequence, err = seq.LoadTextSequence(ac.inputFile)
	} else {
		ac.sequence, err = seq.LoadStructuredSequence(ac.inputFile, ac.format)
	}

	if err != nil {
		return err
	}
	return nil
}

// Comments returns the comments from the loaded sequence
func (ac *AppContext) Comments() []string {
	if ac.sequence == nil {
		return nil
	}
	return ac.sequence.Comments
}

// SampleRate returns the sample rate from the loaded sequence options
func (ac *AppContext) SampleRate() int {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return 0
	}

	return ac.sequence.Options.SampleRate
}

// PresetList returns the preset list from the loaded sequence options
func (ac *AppContext) PresetList() []string {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return []string{}
	}

	return ac.sequence.Options.PresetList
}

// Volume returns the volume from the loaded sequence options
func (ac *AppContext) Volume() int {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return 0
	}

	return ac.sequence.Options.Volume
}

// GainLevel returns the gain level from the loaded sequence options.
// Gain levels:
// 0 = 0 dB,
// 3 = -3 dB,
// 9 = -9 dB,
// 18 = -18 dB
func (ac *AppContext) GainLevel() int {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return 0
	}

	return int(ac.sequence.Options.GainLevel)
}

// BackgroundPath returns the background audio path from the loaded sequence options
func (ac *AppContext) BackgroundPath() string {
	if ac.sequence == nil || ac.sequence.Options == nil {
		return ""
	}

	return ac.sequence.Options.BackgroundPath
}

// RawContent returns the raw content of the loaded sequence
func (ac *AppContext) RawContent() []byte {
	if ac.sequence == nil {
		return nil
	}

	return ac.sequence.RawContent
}
