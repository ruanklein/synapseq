package voices

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// VoiceRegistry holds all registered voice generators
type VoiceRegistry struct {
	generators [t.NumberOfChannels]t.VoiceGenerator
}

// NewVoiceRegistry creates a new voice registry with all built-in generators
func NewVoiceRegistry() *VoiceRegistry {
	registry := &VoiceRegistry{
		generators: [t.NumberOfChannels]t.VoiceGenerator{},
	}

	// Register all built-in generators
	registry.Register(&BinauralGenerator{})
	registry.Register(&MonauralGenerator{})
	// registry.Register(&IsochronicGenerator{})
	// registry.Register(&NoiseGenerator{})
	// registry.Register(&BackgroundGenerator{})
	// registry.Register(&SpinGenerator{})

	return registry
}

// Register adds a voice generator to the registry
func (vr *VoiceRegistry) Register(generator t.VoiceGenerator) {
	voiceType := generator.GetVoiceType()
	if int(voiceType) < len(vr.generators) {
		vr.generators[voiceType] = generator
	}
}

// GetGenerator returns the generator for a specific voice type
func (vr *VoiceRegistry) GetGenerator(voiceType t.VoiceType) t.VoiceGenerator {
	if int(voiceType) < len(vr.generators) {
		return vr.generators[voiceType]
	}
	return nil
}

// GenerateSample generates a sample using the appropriate generator
func (vr *VoiceRegistry) GenerateSample(ch *t.Channel, waveTables [4][]int) (int, int) {
	generator := vr.GetGenerator(ch.Voice.Type)
	if generator == nil {
		return 0, 0 // Silence for unsupported types
	}
	return generator.GenerateSample(ch, waveTables)
}

// UpdateChannel updates channel state using the appropriate generator
func (vr *VoiceRegistry) UpdateChannel(ch *t.Channel, sampleRate int) {
	generator := vr.GetGenerator(ch.Voice.Type)
	if generator != nil {
		generator.UpdateChannel(ch, sampleRate)
	}
}
