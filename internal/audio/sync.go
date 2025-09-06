package audio

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// sync synchronizes the audio renderer state with the current time
func (r *AudioRenderer) sync(timeMs int) {
	// Find the correct period for the current time
	for r.periodIdx+1 < len(r.periods) && timeMs >= r.periods[r.periodIdx+1].Time {
		r.periodIdx++
	}

	if r.periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[r.periodIdx]
	nextTime := timeMs + 1000 // Default next time
	if r.periodIdx+1 < len(r.periods) {
		nextTime = r.periods[r.periodIdx+1].Time
	}

	// Calculate interpolation factor (0.0 to 1.0)
	progress := float64(timeMs-period.Time) / float64(nextTime-period.Time)

	// Update each channel
	for ch := range t.NumberOfChannels {
		if ch >= len(r.channels) || ch >= len(period.VoiceStart) {
			return // Bounds protection
		}

		channel := &r.channels[ch]
		v0 := period.VoiceStart[ch]
		v1 := period.VoiceEnd[ch]

		// Interpolate channel parameters
		channel.Voice.Type = v0.Type
		channel.Voice.Amplitude = t.AmplitudeType(float64(v0.Amplitude)*(1-progress) + float64(v1.Amplitude)*progress)
		channel.Voice.Carrier = v0.Carrier*(1-progress) + v1.Carrier*progress
		channel.Voice.Resonance = v0.Resonance*(1-progress) + v1.Resonance*progress
		channel.Voice.Waveform = v0.Waveform
		channel.Voice.Intensity = t.IntensityType(float64(v0.Intensity)*(1-progress) + float64(v1.Intensity)*progress)

		// Reset offsets if voice type has changed
		if channel.Type != channel.Voice.Type {
			channel.Type = channel.Voice.Type
			channel.Offset[0] = 0
			channel.Offset[1] = 0
		}

		switch channel.Voice.Type {
		case t.VoiceBinauralBeat:
			freq1 := channel.Voice.Carrier + channel.Voice.Resonance/2
			freq2 := channel.Voice.Carrier - channel.Voice.Resonance/2
			channel.Amplitude[0] = int(channel.Voice.Amplitude)
			channel.Amplitude[1] = int(channel.Voice.Amplitude)
			channel.Increment[0] = int(freq1 / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freq2 / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.VoiceMonauralBeat:
			freqHigh := channel.Voice.Carrier + channel.Voice.Resonance/2
			freqLow := channel.Voice.Carrier - channel.Voice.Resonance/2
			channel.Amplitude[0] = int(channel.Voice.Amplitude)
			channel.Increment[0] = int(freqHigh / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freqLow / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.VoiceIsochronicBeat:
			channel.Amplitude[0] = int(channel.Voice.Amplitude)
			channel.Increment[0] = int(channel.Voice.Carrier / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(channel.Voice.Resonance / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.VoiceWhiteNoise, t.VoicePinkNoise, t.VoiceBrownNoise:
			channel.Amplitude[0] = int(channel.Voice.Amplitude)
		case t.VoiceSpinWhite, t.VoiceSpinPink, t.VoiceSpinBrown:
			channel.Amplitude[0] = int(channel.Voice.Amplitude)
			channel.Increment[0] = int(channel.Voice.Resonance / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			// Clamp carrier to a maximum of 127Hz to avoid clicks
			spinCarrierMax := 127.0 / 1e-6 / float64(r.sampleRate)
			clampedCarrier := channel.Voice.Carrier

			if clampedCarrier > spinCarrierMax {
				clampedCarrier = spinCarrierMax
			}
			if clampedCarrier < -spinCarrierMax {
				clampedCarrier = -spinCarrierMax
			}
			channel.Increment[1] = int(clampedCarrier * 1e-6 * float64(r.sampleRate) * float64(1<<24) / float64(t.WaveTableAmplitude))
		}
	}
}
