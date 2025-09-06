package types

import "fmt"

// Period represents a time period with voice configurations
type Period struct {
	Time       int                     // Start time (end time is ->Next->Time)
	VoiceStart [NumberOfChannels]Voice // Start voices for each channel
	VoiceEnd   [NumberOfChannels]Voice // End voices for each channel
}

// TimeString returns the time of this period as a formatted string
func (p *Period) TimeString() string {
	hh := p.Time / 3600000
	mm := (p.Time % 3600000) / 60000
	ss := (p.Time % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

// AdjustPeriods adjusts the voices in the overlapping periods
func AdjustPeriods(last, next *Period) error {
	for ch := range NumberOfChannels {
		v0 := &last.VoiceStart[ch]
		v1 := &last.VoiceEnd[ch]
		v2 := &next.VoiceStart[ch]

		// Apply Fade-In
		if v0.Type == VoiceSilence {
			v0.Type = v2.Type
			v0.Carrier = v2.Carrier
			v0.Resonance = v2.Resonance
			v0.Amplitude = 0
			v0.Intensity = v2.Intensity
			v0.Waveform = v2.Waveform
		}

		// Apply Fade-Out
		if v2.Type == VoiceSilence {
			v2.Carrier = v1.Carrier
			v2.Resonance = v1.Resonance
			v2.Intensity = v1.Intensity
		}

		// Validate if previus period has a voice on and next period turn it off or vice-versa
		if (v1.Type != VoiceOff && v1.Type != VoiceSilence && v2.Type == VoiceOff) ||
			(v1.Type == VoiceOff && v2.Type != VoiceOff && v2.Type != VoiceSilence) {
			return fmt.Errorf("channel %d cannot be turned off or on directly, use silence instead: %s --> %s", ch+1, v1.Type.String(), v2.Type.String())
		}

		// Validate if previus period has a voice on and next period change type
		if (v1.Type != v2.Type) &&
			(v1.Type != VoiceOff &&
				v1.Type != VoiceSilence &&
				v2.Type != VoiceOff &&
				v2.Type != VoiceSilence) {
			return fmt.Errorf("channel %d cannot change voice type directly, use silence instead: %s --> %s", ch+1, v1.Type.String(), v2.Type.String())
		}

		// Validate if previus period has a voice on and next period change waveform
		if (v1.Waveform != v2.Waveform) &&
			(v1.Type != VoiceOff &&
				v1.Type != VoiceSilence &&
				v2.Type != VoiceOff &&
				v2.Type != VoiceSilence) {
			return fmt.Errorf("channel %d cannot change waveform directly, use silence instead: %s --> %s", ch+1, v1.Waveform.String(), v2.Waveform.String())
		}

		// Carry forward the voice settings from the end of the last period to the start of the next period
		v1.Type = v2.Type
		v1.Carrier = v2.Carrier
		v1.Resonance = v2.Resonance
		v1.Amplitude = v2.Amplitude
		v1.Intensity = v2.Intensity
		v1.Waveform = v2.Waveform
	}
	return nil
}
