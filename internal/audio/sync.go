package audio

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// sync synchronizes the audio renderer state with the current time
func (r *AudioRenderer) sync(timeMs int, periodIdx int) {
	if periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[periodIdx]
	nextTime := timeMs + 1000 // Default next time
	if periodIdx+1 < len(r.periods) {
		nextTime = r.periods[periodIdx+1].Time
	}

	// Calculate interpolation factor (0.0 to 1.0)
	progress := float64(timeMs-period.Time) / float64(nextTime-period.Time)

	// Update each channel
	for ch := range t.NumberOfChannels {
		if ch >= len(r.channels) || ch >= len(period.TrackStart) {
			return // Bounds protection
		}

		channel := &r.channels[ch]
		tr0 := period.TrackStart[ch]
		tr1 := period.TrackEnd[ch]

		// Interpolate channel parameters
		channel.Track.Type = tr0.Type
		channel.Track.Effect.Type = tr0.Effect.Type
		channel.Track.Amplitude = t.AmplitudeType(float64(tr0.Amplitude)*(1-progress) + float64(tr1.Amplitude)*progress)
		channel.Track.Carrier = tr0.Carrier*(1-progress) + tr1.Carrier*progress
		channel.Track.Resonance = tr0.Resonance*(1-progress) + tr1.Resonance*progress
		channel.Track.Waveform = tr0.Waveform
		channel.Track.Intensity = t.IntensityType(float64(tr0.Intensity)*(1-progress) + float64(tr1.Intensity)*progress)

		// Reset offsets if track type has changed
		if channel.Type != channel.Track.Type {
			channel.Type = channel.Track.Type
			channel.Offset[0] = 0
			channel.Offset[1] = 0
		}

		switch channel.Track.Type {
		case t.TrackBinauralBeat:
			freq1 := channel.Track.Carrier + channel.Track.Resonance/2
			freq2 := channel.Track.Carrier - channel.Track.Resonance/2
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Amplitude[1] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(freq1 / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freq2 / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackMonauralBeat:
			freqHigh := channel.Track.Carrier + channel.Track.Resonance/2
			freqLow := channel.Track.Carrier - channel.Track.Resonance/2
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(freqHigh / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(freqLow / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackIsochronicBeat:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
			channel.Increment[0] = int(channel.Track.Carrier / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			channel.Increment[1] = int(channel.Track.Resonance / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
		case t.TrackWhiteNoise, t.TrackPinkNoise, t.TrackBrownNoise:
			channel.Amplitude[0] = int(channel.Track.Amplitude)
		case t.TrackBackground:
			channel.Amplitude[0] = int(channel.Track.Amplitude)

			switch channel.Track.Effect.Type {
			case t.EffectSpin:
				channel.Increment[0] = int(channel.Track.Resonance / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)

				spinCarrierMax := 127.0 / 1e-6 / float64(r.sampleRate)
				clampedCarrier := channel.Track.Carrier

				if clampedCarrier > spinCarrierMax {
					clampedCarrier = spinCarrierMax
				}
				if clampedCarrier < -spinCarrierMax {
					clampedCarrier = -spinCarrierMax
				}
				channel.Increment[1] = int(clampedCarrier * 1e-6 * float64(r.sampleRate) * float64(1<<24) / float64(t.WaveTableAmplitude))
			case t.EffectPulse:
				channel.Increment[1] = int(channel.Track.Resonance / float64(r.sampleRate) * t.SineTableSize * t.PhasePrecision)
			}
		}
	}
}
