/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

const (
	NumberOfChannels = 16 // Max number of channels
)

// Channel represents a channel state
type Channel struct {
	// Current track setting (updated from current period)
	Track Track
	// Track type
	Type TrackType
	// Current amplitude state
	Amplitude [2]int
	// Increment (for binaural tones, offset + increment into sine table * 65536)
	Increment [2]int
	// Offset into waveform table (for tones, offset + increment into sine table * 65536)
	Offset [2]int
}
