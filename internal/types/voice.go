package types

import "fmt"

// VoiceType represents the type of voice/sound
type VoiceType int

// Voice types
const (
	VoiceOff            VoiceType = iota // Voice is off
	VoiceSilence                         // Voice is silence
	VoiceBinauralBeat                    // Voice is a binaural beat
	VoiceMonauralBeat                    // Voice is a monaural beat
	VoiceIsochronicBeat                  // Voice is an isochronic beat
	VoiceWhiteNoise                      // Voice is white noise
	VoicePinkNoise                       // Voice is pink noise
	VoiceBrownNoise                      // Voice is brown noise
	VoiceSpinWhite                       // Voice is a spin white noise
	VoiceSpinPink                        // Voice is a spin pink noise
	VoiceSpinBrown                       // Voice is a spin brown noise
	VoiceBackground                      // Voice is a background noise
	VoiceEffectSpin                      // Voice is a spin effect
	VoiceEffectPulse                     // Voice is a pulse effect
)

// String returns the string representation of the VoiceType
func (vt VoiceType) String() string {
	switch vt {
	case VoiceBinauralBeat:
		return KeywordBinaural
	case VoiceMonauralBeat:
		return KeywordMonaural
	case VoiceIsochronicBeat:
		return KeywordIsochronic
	case VoiceWhiteNoise, VoiceSpinWhite:
		return KeywordWhite
	case VoicePinkNoise, VoiceSpinPink:
		return KeywordPink
	case VoiceBrownNoise, VoiceSpinBrown:
		return KeywordBrown
	case VoiceBackground:
		return KeywordBackground
	case VoiceEffectSpin:
		return KeywordSpin
	case VoiceEffectPulse:
		return KeywordPulse
	default:
		return "- -"
	}
}

// Voice represents a voice configuration
type Voice struct {
	Type      VoiceType     // Voice type
	Amplitude AmplitudeType // Amplitude level (0-4096 for 0-100%)
	Carrier   float64       // Carrier frequency
	Resonance float64       // Resonance frequency
	Waveform  WaveformType  // Waveform shape
	Intensity IntensityType // Intensity (for effects)
}

// Validate checks if the voice configuration is valid
func (v *Voice) Validate() error {
	if v.Amplitude < 0 || v.Amplitude > 4096 {
		return fmt.Errorf("amplitude must be between 0 and 100. Received: %.2f", v.Amplitude.ToPercent())
	}
	if v.Carrier < 0 {
		return fmt.Errorf("carrier frequency must be positive: %.2f. Received: %.2f", v.Carrier, v.Carrier)
	}
	if v.Resonance < 0 {
		return fmt.Errorf("resonance frequency must be positive: %.2f. Received: %.2f", v.Resonance, v.Resonance)
	}
	if v.Intensity < 0 || v.Intensity > 1.0 {
		return fmt.Errorf("intensity must be between 0 and 100. Received: %.2f", v.Intensity.ToPercent())
	}
	return nil
}

// String returns the string representation of the Voice
func (v *Voice) String() string {
	switch v.Type {
	case VoiceOff:
		return "- -"
	case VoiceBinauralBeat, VoiceMonauralBeat, VoiceIsochronicBeat:
		return fmt.Sprintf("%s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, v.Waveform.String(), KeywordTone, v.Carrier, v.Type.String(), v.Resonance, KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceWhiteNoise, VoicePinkNoise, VoiceBrownNoise:
		return fmt.Sprintf("%s %s %s %.2f", KeywordNoise, v.Type.String(), KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceSpinWhite, VoiceSpinPink, VoiceSpinBrown:
		return fmt.Sprintf("%s %s %s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, v.Waveform.String(), KeywordSpin, v.Type.String(), KeywordWidth, v.Carrier, KeywordRate, v.Resonance, KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceBackground:
		return fmt.Sprintf("%s %s %.2f", KeywordBackground, KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceEffectSpin:
		return fmt.Sprintf("%s %s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, v.Waveform.String(), KeywordSpin, KeywordWidth, v.Carrier, KeywordRate, v.Resonance, KeywordIntensity, v.Intensity.ToPercent())
	case VoiceEffectPulse:
		return fmt.Sprintf("%s %s %s %.2f %s %.2f", KeywordWaveform, v.Waveform.String(), KeywordPulse, v.Resonance, KeywordIntensity, v.Intensity.ToPercent())
	default:
		return ""
	}
}

// CompactString returns a compact string representation of the voice configuration
func (v *Voice) CompactString() string {
	switch v.Type {
	case VoiceOff:
		return " -"
	case VoiceBinauralBeat, VoiceMonauralBeat, VoiceIsochronicBeat:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordTone, v.Carrier, v.Type.String(), v.Resonance, KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceWhiteNoise, VoicePinkNoise, VoiceBrownNoise:
		return fmt.Sprintf(" (%s:%.2f)", KeywordNoise, v.Amplitude.ToPercent())
	case VoiceSpinWhite, VoiceSpinPink, VoiceSpinBrown:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordSpin, v.Carrier, KeywordRate, v.Resonance, KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceBackground:
		return fmt.Sprintf(" (%s:%.2f)", KeywordAmplitude, v.Amplitude.ToPercent())
	case VoiceEffectSpin:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordSpin, v.Carrier, KeywordRate, v.Resonance, KeywordIntensity, v.Intensity.ToPercent())
	case VoiceEffectPulse:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f)", KeywordPulse, v.Resonance, KeywordIntensity, v.Intensity.ToPercent())
	default:
		return " ???"
	}
}

// Equal checks if two voice configurations are identical
func (v *Voice) Equal(v2 *Voice) bool {
	return v.Type == v2.Type &&
		v.Amplitude == v2.Amplitude &&
		v.Carrier == v2.Carrier &&
		v.Resonance == v2.Resonance &&
		v.Waveform == v2.Waveform &&
		v.Intensity == v2.Intensity
}
