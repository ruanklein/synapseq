package sequence

import (
	"fmt"

	"github.com/ruanklein/synapseq/internal/audio"
)

const (
	MaxPresets = 32 // Maximum number of presets

	builtinSilence = "silence" // Represents silence built-in preset
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
		p.Voice[i].Intensity = 0.0
	}
}

// AllocateVoice allocates a free voice in the preset
func (p *Preset) AllocateVoice() (int, error) {
	for index, voice := range p.Voice {
		if voice.Type == audio.VoiceOff {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no available voices for preset '%s'", p.Name)
}

// AllVoicesAreOff checks if all voices in the preset are off
func (p *Preset) AllVoicesAreOff() bool {
	for _, voice := range p.Voice {
		if voice.Type != audio.VoiceOff {
			return false
		}
	}
	return true
}

// GetBackgroundVoice retrieves the background voice if it exists
func (p *Preset) GetBackgroundVoice() *audio.Voice {
	for i := range p.Voice {
		if p.Voice[i].Type == audio.VoiceBackground {
			return &p.Voice[i]
		}
	}
	return nil
}

func (p *Preset) HasNext() bool {
	return p.Next != nil
}
