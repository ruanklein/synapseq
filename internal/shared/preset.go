package shared

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// FindPreset searches for a preset by name in a slice of presets
func FindPreset(name string, presets []t.Preset) *t.Preset {
	for i := range presets {
		if presets[i].String() == name {
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
	return -1, fmt.Errorf("no available voices for preset %q", preset.String())
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

// NumBackgroundVoices counts the number of background voices in the preset
func NumBackgroundVoices(preset *t.Preset) int {
	count := 0
	for _, voice := range preset.Voice {
		if voice.Type == t.VoiceBackground {
			count++
		}
	}
	return count
}
