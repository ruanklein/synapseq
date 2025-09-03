package generator

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// GenerateSample generates a sample using the appropriate generator
func GenerateSample(ch *t.Channel, waveTables [4][]int) (int, int) {
	switch ch.Voice.Type {
	case t.VoiceBinauralBeat:
		return binauralGenerateSample(ch, waveTables)
	case t.VoiceMonauralBeat:
		return monauralGenerateSample(ch, waveTables)
	default:
		return 0, 0 // Silence for unsupported types
	}
}

// UpdateChannel updates channel state using the appropriate generator
func UpdateChannel(ch *t.Channel, sampleRate int) {
	switch ch.Voice.Type {
	case t.VoiceBinauralBeat:
		binauralUpdateChannel(ch, sampleRate)
	case t.VoiceMonauralBeat:
		monauralUpdateChannel(ch, sampleRate)
	default:
		// Silence for unsupported types
	}
}
