/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// calcPulseFactor calculates the pulse effect modulation factor for a channel
func (r *AudioRenderer) calcPulseFactor(channel *t.Channel) float64 {
	modVal := float64(r.waveTables[int(channel.Track.Waveform)][channel.Offset[1]>>16])

	threshold := 0.3 * float64(t.WaveTableAmplitude)
	den := 0.7 * float64(t.WaveTableAmplitude)

	modFactor := 0.0
	if modVal > threshold {
		modFactor = (modVal - threshold) / den
		modFactor = modFactor * modFactor * (3 - 2*modFactor)
	}

	return modFactor
}
