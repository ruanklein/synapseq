/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"math"

	t "github.com/ruanklein/synapseq/internal/types"
)

// mix generates a stereo audio sample by mixing all channels
func (r *AudioRenderer) mix(samples []int) []int {
	// Function to get the next dither value
	nextDither := func() int64 {
		r.dither0 = r.dither1
		r.dither1 = uint16((uint32(r.dither0)*0x660D + 0xF35F) & 0xFFFF)
		return int64(int32(r.dither1) - 32768) // ~[-32768..32767]
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
			case t.TrackPureTone:
				channel.Offset[0] += channel.Increment[0]
				channel.Offset[0] &= (t.SineTableSize << 16) - 1

				left += int64(channel.Amplitude[0]) * int64(r.waveTables[waveIdx][channel.Offset[0]>>16])
				right += int64(channel.Amplitude[0]) * int64(r.waveTables[waveIdx][channel.Offset[0]>>16])
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
				// Use pre-generated pink noise sample for efficiency
				noiseVal := int64(r.noiseGenerator.Generate(t.TrackPinkNoise))
				if channel.Track.Type != t.TrackPinkNoise {
					noiseVal = int64(r.noiseGenerator.Generate(channel.Track.Type))
				}

				// Scale noise by amplitude
				sampleVal := int64(channel.Amplitude[0]) * noiseVal

				left += sampleVal
				right += sampleVal
			case t.TrackBackground:
				// Gain factor from dB to linear (0.0..1.0)
				// -20dB = 0.1, -12dB ≈ 0.251, -6dB ≈ 0.5, 0dB = 1.0
				dbValue := -float64(r.GainLevel)
				g := math.Pow(10, dbValue/20.0)

				// Background audio sample (stereo interleaved)
				bgLFloat := float64(backgroundSamples[i*2]) * g
				bgRFloat := float64(backgroundSamples[i*2+1]) * g

				// Amplitude scaling (0..256)
				// Divide by 16 to convert 0..4096 to 0..256
				// This allows finer control over background volume
				// without exceeding int64 limits during mixing
				bgAmp := float64(channel.Amplitude[0]) / 16.0 // 0..256

				// Final background sample values
				bgL := int64(bgLFloat * bgAmp)
				bgR := int64(bgRFloat * bgAmp)

				switch channel.Track.Effect.Type {
				case t.EffectSpin:
					channel.Offset[0] += channel.Increment[0]
					channel.Offset[0] &= (t.SineTableSize << 16) - 1

					spinPos := (channel.Increment[1] * r.waveTables[waveIdx][channel.Offset[0]>>16]) >> 24

					effectIntensity := float64(channel.Track.Intensity) * 0.7
					spinGain := 0.5 + effectIntensity*3.5

					ampSpin := int64(float64(spinPos) * spinGain)
					if ampSpin > 127 {
						ampSpin = 127
					}
					if ampSpin < -128 {
						ampSpin = -128
					}

					posVal := ampSpin
					if posVal < 0 {
						posVal = -posVal
					}
					if posVal > 128 {
						posVal = 128
					}

					var spinLeft, spinRight int64
					if ampSpin >= 0 {
						// L = BG_L * (128 - pos)/128
						// R = BG_R + BG_L * pos/128
						spinLeft = (bgL * (128 - posVal)) >> 7
						spinRight = bgR + ((bgL * posVal) >> 7)
					} else {
						// L = BG_L + BG_R * pos/128
						// R = BG_R * (128 - pos)/128
						spinLeft = bgL + ((bgR * posVal) >> 7)
						spinRight = (bgR * (128 - posVal)) >> 7
					}

					left += spinLeft
					right += spinRight

				case t.EffectPulse:
					// LFO for pulse modulation
					channel.Offset[1] += channel.Increment[1]
					channel.Offset[1] &= (t.SineTableSize << 16) - 1

					// 0..1
					modFactor := r.calcPulseFactor(channel)

					// Mix the effect (0..1) weighted by intensity
					effectIntensity := float64(channel.Track.Intensity) * 0.7
					gain := (1.0 - effectIntensity) + (effectIntensity * modFactor)

					left += int64(float64(bgL) * gain)
					right += int64(float64(bgR) * gain)

				default:
					// BG without effect
					left += bgL
					right += bgR
				}
			}
		}

		if r.Volume != 100 {
			left = left * int64(r.Volume) / 100
			right = right * int64(r.Volume) / 100
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
