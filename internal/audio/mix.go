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

	// Read background audio samples if enabled
	var backgroundSamples []int

	if r.backgroundAudio.IsEnabled() {
		// Buffer for background audio
		backgroundSamples = make([]int, t.BufferSize*audioChannels) // Stereo
		r.backgroundAudio.ReadSamples(backgroundSamples, t.BufferSize*audioChannels)
	}

	for i := range t.BufferSize {
		var left, right int64

		for ch := range t.NumberOfChannels {
			channel := &r.channels[ch]
			waveIdx := int(channel.Track.Waveform)

			switch channel.Track.Type {
			case t.TrackBinauralBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				left += int64(channel.Amplitude[0]) * int64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				right += int64(channel.Amplitude[1]) * int64(r.waveTables[waveIdx][channel.Offset[1]>>16])
			case t.TrackMonauralBeat:
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
			case t.TrackIsochronicBeat:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				channel.Offset[1] += channel.Increment[1]
				channel.Offset[1] &= (t.SineTableSize << 16) - 1

				modFactor := r.calcPulseFactor(channel)

				carrier := float64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				amp := float64(channel.Amplitude[0])

				out := int64(amp * carrier * modFactor)

				left += out
				right += out
			case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
				noiseVal := int64(r.noiseGenerator.Generate(channel.Track.Type))
				sampleVal := int64(channel.Amplitude[0]) * noiseVal

				left += sampleVal
				right += sampleVal
			case t.TrackBackground:
				g := float64(r.gainLevel) / 100.0
				bgLeft := int64(float64(backgroundSamples[i*2]) * g)
				bgRight := int64(float64(backgroundSamples[i*2+1]) * g)

				backgroundAmplitude := int64(channel.Amplitude[0])

				switch channel.Track.Effect.Type {
				case t.EffectSpin:
					channel.Offset[0] += channel.Increment[0]
					channel.Offset[0] &= (t.SineTableSize << 16) - 1

					spinPos := (channel.Increment[1] * r.waveTables[waveIdx][channel.Offset[0]>>16]) >> 24

					effectIntensity := float64(channel.Track.Intensity) * 0.7
					amplifiedSpin := int64(float64(spinPos) * (0.5 + effectIntensity*3.5))

					if amplifiedSpin > 127 {
						amplifiedSpin = 127
					}
					if amplifiedSpin < -128 {
						amplifiedSpin = -128
					}

					posVal := amplifiedSpin
					if posVal < 0 {
						posVal = -posVal
					}

					var spinLeft, spinRight int64
					if amplifiedSpin >= 0 {
						spinLeft = (bgLeft * backgroundAmplitude * (128 - posVal)) >> 7
						spinRight = bgRight*backgroundAmplitude + ((bgLeft * backgroundAmplitude * posVal) >> 7)
					} else {
						spinLeft = bgLeft*backgroundAmplitude + ((bgRight * backgroundAmplitude * posVal) >> 7)
						spinRight = (bgRight * backgroundAmplitude * (128 - posVal)) >> 7
					}

					left += spinLeft
					right += spinRight
				case t.EffectPulse:
					channel.Offset[1] += channel.Increment[1]
					channel.Offset[1] &= (t.SineTableSize << 16) - 1

					modFactor := r.calcPulseFactor(channel)

					effectIntensity := float64(channel.Track.Intensity) * 0.7
					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)

					left += int64(float64(bgLeft*backgroundAmplitude) * gain)
					right += int64(float64(bgRight*backgroundAmplitude) * gain)
				default:
					left += bgLeft * backgroundAmplitude
					right += bgRight * backgroundAmplitude
				}
			}
		}

		if r.volume != 100 {
			left = left * int64(r.volume) / 100
			right = right * int64(r.volume) / 100
		}

		// Apply dithering
		left += nextDither()
		right += nextDither()

		// Scale down to 24-bit range
		left >>= audioBitShift
		right >>= audioBitShift

		// Clipping to 24-bit range
		if left > audioMaxValue {
			left = audioMaxValue
		}
		if left < audioMinValue {
			left = audioMinValue
		}
		if right > audioMaxValue {
			right = audioMaxValue
		}
		if right < audioMinValue {
			right = audioMinValue
		}

		samples[i*2] = int(left)
		samples[i*2+1] = int(right)
	}

	return samples
}
