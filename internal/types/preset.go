package types

import (
	"fmt"
	"strings"
)

const (
	MaxPresets     = 32        // Maximum number of presets
	builtinSilence = "silence" // Represents silence built-in preset
)

// Preset represents a named preset
type Preset struct {
	name  string                  // Name of preset
	Voice [NumberOfChannels]Voice // Voice-set for it
}

// NewPreset creates a new preset with the given name
func NewPreset(name string) (*Preset, error) {
	isLetter := func(b byte) bool {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
	}
	isDigit := func(b byte) bool {
		return b >= '0' && b <= '9'
	}

	if len(name) == 0 {
		return nil, fmt.Errorf("preset name cannot be empty")
	}

	first := name[0]
	if !isLetter(first) {
		return nil, fmt.Errorf("preset name must start with a letter: %q", name)
	}

	for i := 1; i < len(name); i++ {
		ch := name[i]
		if !(isLetter(ch) || isDigit(ch) || ch == '_' || ch == '-') {
			return nil, fmt.Errorf("invalid character in preset name %q: %q", name, string(ch))
		}
	}

	n := strings.ToLower(name)
	if n == builtinSilence {
		return nil, fmt.Errorf("preset name %q is reserved", builtinSilence)
	}

	preset := &Preset{name: n}
	for i := range NumberOfChannels {
		preset.Voice[i].Type = VoiceOff
		preset.Voice[i].Carrier = 0.0
		preset.Voice[i].Resonance = 0.0
		preset.Voice[i].Amplitude = 0.0
		preset.Voice[i].Intensity = 0.0
		preset.Voice[i].Waveform = WaveformSine
	}
	return preset, nil
}

// NewBuiltinSilencePreset creates a new silence preset
func NewBuiltinSilencePreset() *Preset {
	preset := &Preset{name: builtinSilence}
	for i := range NumberOfChannels {
		preset.Voice[i].Type = VoiceSilence
		preset.Voice[i].Carrier = 0.0
		preset.Voice[i].Resonance = 0.0
		preset.Voice[i].Amplitude = 0.0
		preset.Voice[i].Intensity = 0.0
		preset.Voice[i].Waveform = WaveformSine
	}
	return preset
}

// String returns the name of the preset
func (p *Preset) String() string {
	return p.name
}
