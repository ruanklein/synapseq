package audio

// VoiceType represents the type of voice/sound
type VoiceType int

// Voice types
const (
	VoiceOff            VoiceType = iota // Voice is off
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

// Voice represents a voice configuration
type Voice struct {
	Type      VoiceType    // Voice type
	Amplitude float64      // Amplitude level (0-4096 for 0-100%)
	Carrier   float64      // Carrier frequency
	Resonance float64      // Resonance frequency
	Waveform  WaveformType // Waveform shape
}

// Equal checks if two Voice structs are equal
func (v1 Voice) Equal(v2 Voice) bool {
	return v1.Type == v2.Type &&
		v1.Amplitude == v2.Amplitude &&
		v1.Carrier == v2.Carrier &&
		v1.Resonance == v2.Resonance &&
		v1.Waveform == v2.Waveform
}

// IsOff checks if the voice is off
func (v Voice) IsOff() bool {
	return v.Type == VoiceOff || v.Amplitude == 0
}
