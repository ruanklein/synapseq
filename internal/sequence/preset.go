package sequence

import (
	"regexp"

	"github.com/ruanklein/synapseq/internal/audio"
)

const (
	BuiltinSilence = "silence" // Represents silence
)

var (
	PresetList []Preset // List of presets
)

// Preset represents a named preset
type Preset struct {
	Next  *Preset                             // Next preset in list
	Name  string                              // Name of preset
	Voice [audio.NumberOfChannels]audio.Voice // Voice-set for it
}

// InitVoices initializes the voices in the preset
func (p *Preset) InitVoices() {
	for i := range audio.NumberOfChannels {
		p.Voice[i].Type = audio.VoiceOff
		p.Voice[i].Amplitude = 0.0
		p.Voice[i].Carrier = 0.0
		p.Voice[i].Resonance = 0.0
		p.Voice[i].Waveform = audio.WaveformSine
	}
}

// Regex for validating preset names
func IsPreset(s string) bool {
	regexPreset := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	return regexPreset.MatchString(s)
}
