package audio

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// mix generates a stereo audio sample by mixing all channels
func (r *AudioRenderer) mix(samples []int) []int {
	for i := range t.BufferSize {
		left, right := 0, 0

		for ch := range t.NumberOfChannels {
			channel := &r.channels[ch]
			waveIdx := int(channel.Voice.Waveform)

			switch channel.Voice.Type {
			case t.VoiceBinauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				left += channel.Amplitude[0] * r.waveTables[waveIdx][channel.Offset[0]>>16]
				right += channel.Amplitude[1] * r.waveTables[waveIdx][channel.Offset[1]>>16]
			case t.VoiceMonauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				freqHigh := r.waveTables[waveIdx][channel.Offset[0]>>16]
				freqLow := r.waveTables[waveIdx][channel.Offset[1]>>16]

				halfAmp := channel.Amplitude[0] / 2
				mixedSample := halfAmp * (freqHigh + freqLow)

				left += mixedSample
				right += mixedSample
			case t.VoiceIsochronicBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				modVal := float64(r.waveTables[waveIdx][channel.Offset[1]>>16])
				threshold := 0.3 * float64(t.WaveTableAmplitude)
				den := 0.7 * float64(t.WaveTableAmplitude)

				factor := 0.0
				if modVal > threshold {
					factor = (modVal - threshold) / den
					factor = factor * factor * (3 - 2*factor)
				}

				carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				amp := float64(channel.Amplitude[0])

				out := int(amp * carrier * factor)

				left += out
				right += out
			}
		}

		if r.volume != 100 {
			left = int(int64(left) * int64(r.volume) / 100)
			right = int(int64(right) * int64(r.volume) / 100)
		}

		left >>= 16
		right >>= 16

		// Clipping to 16-bit range
		if left > 32767 {
			left = 32767
		}
		if left < -32768 {
			left = -32768
		}
		if right > 32767 {
			right = 32767
		}
		if right < -32768 {
			right = -32768
		}

		samples[i*2] = left
		samples[i*2+1] = right
	}

	return samples
}
