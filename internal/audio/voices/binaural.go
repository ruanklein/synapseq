package voices

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// BinauralGenerator generates binaural beat samples
type BinauralGenerator struct{}

// GetVoiceType returns the voice type this generator handles
func (bg *BinauralGenerator) GetVoiceType() t.VoiceType {
	return t.VoiceBinauralBeat
}

// UpdateChannel updates channel state for binaural beats
func (bg *BinauralGenerator) UpdateChannel(ch *t.Channel, sampleRate int) {
	freq1 := ch.Voice.Carrier + ch.Voice.Resonance/2
	freq2 := ch.Voice.Carrier - ch.Voice.Resonance/2
	ch.Amplitude[0] = int(ch.Voice.Amplitude)
	ch.Amplitude[1] = int(ch.Voice.Amplitude)
	ch.Increment[0] = int(freq1 / float64(sampleRate) * t.SineTableSize * 65536)
	ch.Increment[1] = int(freq2 / float64(sampleRate) * t.SineTableSize * 65536)
}

// GenerateSample generates a binaural beat sample
func (bg *BinauralGenerator) GenerateSample(ch *t.Channel, waveTables [4][]int) (int, int) {
	// Advance offset for each ear
	ch.Offset[0] += ch.Increment[0]
	ch.Offset[0] &= (t.SineTableSize << 16) - 1

	ch.Offset[1] += ch.Increment[1]
	ch.Offset[1] &= (t.SineTableSize << 16) - 1

	// Generate samples using waveform table
	waveIdx := int(ch.Voice.Waveform) % 4
	if waveIdx >= len(waveTables) {
		waveIdx = 0 // Default to sine wave
	}

	leftSample := ch.Amplitude[0] * waveTables[waveIdx][ch.Offset[0]>>16]
	rightSample := ch.Amplitude[1] * waveTables[waveIdx][ch.Offset[1]>>16]

	return leftSample, rightSample
}
