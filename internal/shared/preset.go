package shared

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// FindPreset searches for a preset by name in a slice of presets
func FindPreset(name string, presets []t.Preset) *t.Preset {
	for i := range presets {
		if presets[i].Name == name {
			return &presets[i]
		}
	}
	return nil
}

// AllocateVoice allocates a free voice in the preset
func AllocateVoice(preset *t.Preset) (int, error) {
	for index, voice := range preset.Voice {
		if voice.Type == t.VoiceOff {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no available voices for preset %q", preset.Name)
}

// IsPresetEmpty checks if all voices in the preset are off
func IsPresetEmpty(preset *t.Preset) bool {
	for _, voice := range preset.Voice {
		if voice.Type != t.VoiceOff {
			return false
		}
	}
	return true
}

// InitPresetVoices initializes the voices in the preset
func InitPresetVoices(preset *t.Preset, defaultType t.VoiceType) {
	for i := range t.NumberOfChannels {
		preset.Voice[i].Type = defaultType
		preset.Voice[i].Amplitude = 0.0
		preset.Voice[i].Carrier = 0.0
		preset.Voice[i].Resonance = 0.0
		preset.Voice[i].Waveform = t.WaveformSine
		preset.Voice[i].Intensity = 0.0
	}
}
