package shared

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// AdjustPeriods adjusts the tracks in the overlapping periods
func AdjustPeriods(last, next *t.Period) error {
	for ch := range t.NumberOfChannels {
		tr0 := &last.TrackStart[ch]
		tr1 := &last.TrackEnd[ch]
		tr2 := &next.TrackStart[ch]

		// Apply Fade-In
		if tr0.Type == t.TrackSilence {
			tr0.Type = tr2.Type
			tr0.Carrier = tr2.Carrier
			tr0.Resonance = tr2.Resonance
			tr0.Amplitude = 0
			tr0.Intensity = tr2.Intensity
			tr0.Waveform = tr2.Waveform
		}

		// Apply Fade-Out
		if tr2.Type == t.TrackSilence {
			tr2.Carrier = tr1.Carrier
			tr2.Resonance = tr1.Resonance
			tr2.Intensity = tr1.Intensity
		}

		// Validate if previus period has a track on and next period turn it off or vice-versa
		if (tr1.Type != t.TrackOff && tr1.Type != t.TrackSilence && tr2.Type == t.TrackOff) ||
			(tr1.Type == t.TrackOff && tr2.Type != t.TrackOff && tr2.Type != t.TrackSilence) {
			return fmt.Errorf("channel %d cannot be turned off or on directly, use silence instead: %s --> %s", ch+1, tr1.Type.String(), tr2.Type.String())
		}

		// Validate if previus period has a track on and next period change type
		if (tr1.Type != tr2.Type) &&
			(tr1.Type != t.TrackOff &&
				tr1.Type != t.TrackSilence &&
				tr2.Type != t.TrackOff &&
				tr2.Type != t.TrackSilence) {
			return fmt.Errorf("channel %d cannot change track type directly, use silence instead: %s --> %s", ch+1, tr1.Type.String(), tr2.Type.String())
		}

		// Validate if previus period has a track on and next period change waveform
		if (tr1.Waveform != tr2.Waveform) &&
			(tr1.Type != t.TrackOff &&
				tr1.Type != t.TrackSilence &&
				tr2.Type != t.TrackOff &&
				tr2.Type != t.TrackSilence) {
			return fmt.Errorf("channel %d cannot change waveform directly, use silence instead: %s --> %s", ch+1, tr1.Waveform.String(), tr2.Waveform.String())
		}

		// Carry forward the track settings from the end of the last period to the start of the next period
		tr1.Type = tr2.Type
		tr1.Carrier = tr2.Carrier
		tr1.Resonance = tr2.Resonance
		tr1.Amplitude = tr2.Amplitude
		tr1.Intensity = tr2.Intensity
		tr1.Waveform = tr2.Waveform
	}
	return nil
}
