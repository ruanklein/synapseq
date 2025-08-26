package types

import (
	"fmt"
)

const (
	MaxPresets     = 32        // Maximum number of presets
	BuiltinSilence = "silence" // Represents silence built-in preset
)

// Preset represents a named preset
type Preset struct {
	Next  *Preset                 // Next preset in list
	Name  string                  // Name of preset
	Voice [NumberOfChannels]Voice // Voice-set for it
}

// InitVoices initializes the voices in the preset
func (p *Preset) InitVoices() {
	for i := range NumberOfChannels {
		p.Voice[i].Type = VoiceOff
		p.Voice[i].Amplitude = 0.0
		p.Voice[i].Carrier = 0.0
		p.Voice[i].Resonance = 0.0
		p.Voice[i].Waveform = WaveformSine
		p.Voice[i].Intensity = 0.0
	}
}

// AllocateVoice allocates a free voice in the preset
func (p *Preset) AllocateVoice() (int, error) {
	for index, voice := range p.Voice {
		if voice.Type == VoiceOff {
			return index, nil
		}
	}
	return -1, fmt.Errorf("no available voices for preset '%s'", p.Name)
}

// AllVoicesAreOff checks if all voices in the preset are off
func (p *Preset) AllVoicesAreOff() bool {
	for _, voice := range p.Voice {
		if voice.Type != VoiceOff {
			return false
		}
	}
	return true
}

// GetBackgroundVoice retrieves the background voice if it exists
func (p *Preset) GetBackgroundVoice() *Voice {
	for i := range p.Voice {
		if p.Voice[i].Type == VoiceBackground {
			return &p.Voice[i]
		}
	}
	return nil
}

// HasNext checks if the preset has a next preset
func (p *Preset) HasNext() bool {
	return p.Next != nil
}
