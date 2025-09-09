package audio

import (
	t "github.com/ruanklein/synapseq/internal/types"
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
