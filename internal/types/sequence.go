/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

import "fmt"

// Sequence represents a brainwave sequence
type Sequence struct {
	Periods  []Period
	Options  *SequenceOptions
	Comments []string
}

// SequenceOptions represents configuration options for a sequence
type SequenceOptions struct {
	// Sample rate (e.g., 44100)
	SampleRate int
	// Volume level (0-100 for 0-100%)
	Volume int
	// Path to the background audio file
	BackgroundPath string
	// Path to the preset configuration file
	PresetPath string
	// Gain level (20, 16, 12, 6, 0) for audio processing
	GainLevel GainLevel
}

// Validate checks if the sequence options are valid
func (so *SequenceOptions) Validate() error {
	if so.SampleRate <= 0 {
		return fmt.Errorf("invalid sample rate: %d", so.SampleRate)
	}
	if so.Volume < 0 || so.Volume > 100 {
		return fmt.Errorf("invalid volume: %d", so.Volume)
	}
	return nil
}
