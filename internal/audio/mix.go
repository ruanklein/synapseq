package audio

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// mix generates a stereo audio sample by mixing all channels
func (r *AudioRenderer) mix(samples []int) []int {
	// Simple linear congruential generator for dithering
	var ditherState uint32 = 0x12345678

	// Function to get the next dither value
	nextDither := func() int64 {
		ditherState = ditherState*1103515245 + 12345
		return int64(int32(ditherState>>16) - 32768)
	}

	for i := range t.BufferSize {
		var left, right int64

		for ch := range t.NumberOfChannels {
			channel := &r.channels[ch]
			waveIdx := int(channel.Voice.Waveform)

			switch channel.Voice.Type {
			case t.VoiceBinauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				left += int64(channel.Amplitude[0]) * int64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				right += int64(channel.Amplitude[1]) * int64(r.waveTables[waveIdx][channel.Offset[1]>>16])
			case t.VoiceMonauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				freqHigh := int64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				freqLow := int64(r.waveTables[waveIdx][channel.Offset[1]>>16])

				halfAmp := int64(channel.Amplitude[0]) / 2
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

				out := int64(amp * carrier * factor)

				left += out
				right += out
			case t.VoiceWhiteNoise, t.VoicePinkNoise, t.VoiceBrownNoise:
				noiseVal := int64(r.noiseGenerator.Generate(channel.Voice.Type))
				sampleVal := int64(channel.Amplitude[0]) * noiseVal

				left += sampleVal
				right += sampleVal
			case t.VoiceSpinWhite, t.VoiceSpinPink, t.VoiceSpinBrown:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				spinPos := (channel.Increment[1] * r.waveTables[waveIdx][channel.Offset[0]>>16]) >> 24
				spinLeft, spinRight := r.noiseGenerator.GenerateSpinEffect(channel.Voice.Type, channel.Amplitude[0], spinPos)

				left += spinLeft
				right += spinRight
			}
		}

		if r.volume != 100 {
			left = left * int64(r.volume) / 100
			right = right * int64(r.volume) / 100
		}

		// Apply dithering
		left += nextDither()
		right += nextDither()

		// Scale down to 16-bit range
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

		samples[i*2] = int(left)
		samples[i*2+1] = int(right)
	}

	return samples
}
