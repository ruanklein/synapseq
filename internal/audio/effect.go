/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package audio

import (
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
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
