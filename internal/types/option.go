/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

import (
	"fmt"
	"os"
)

// Option represents configuration options for audio processing.
type Option struct {
	SampleRate     int       // Sample rate (e.g., 44100)
	Volume         int       // Volume level (0-100 for 0-100%)
	BackgroundPath string    // Path to the background audio file
	PresetPath     string    // Path to the preset configuration file
	GainLevel      GainLevel // Gain level (20, 16, 12, 6, 0) for audio processing
}

// Validate checks if the options are valid
func (o *Option) Validate() error {
	if o.SampleRate <= 0 {
		return fmt.Errorf("invalid sample rate: %d", o.SampleRate)
	}
	if o.Volume < 0 || o.Volume > 100 {
		return fmt.Errorf("invalid volume: %d", o.Volume)
	}
	if o.BackgroundPath != "" {
		if _, err := os.Stat(o.BackgroundPath); os.IsNotExist(err) {
			return fmt.Errorf("background path does not exist: %s", o.BackgroundPath)
		}
	}
	if o.PresetPath != "" {
		if _, err := os.Stat(o.PresetPath); os.IsNotExist(err) {
			return fmt.Errorf("preset path does not exist: %s", o.PresetPath)
		}
	}
	return nil
}
