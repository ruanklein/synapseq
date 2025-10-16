/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package shared

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// FindPreset searches for a preset by name in a slice of presets
func FindPreset(name string, presets []t.Preset) *t.Preset {
	for _, preset := range presets {
		if preset.String() == name {
			return &preset
		}
	}
	return nil
}

// AllocateTrack allocates a free track in the preset
func AllocateTrack(preset *t.Preset) (int, error) {
	for index, track := range preset.Track {
		if track.Type == t.TrackOff {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no available tracks for preset %q", preset.String())
}

// IsPresetEmpty checks if all tracks in the preset are off
func IsPresetEmpty(preset *t.Preset) bool {
	for _, track := range preset.Track {
		if track.Type != t.TrackOff {
			return false
		}
	}
	return true
}

// NumBackgroundTracks counts the number of background tracks in the preset
func NumBackgroundTracks(preset *t.Preset) int {
	count := 0
	for _, track := range preset.Track {
		if track.Type == t.TrackBackground {
			count++
		}
	}
	return count
}
