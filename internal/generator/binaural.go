package generator

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// binauralUpdateChannel updates channel state for binaural beats
func binauralUpdateChannel(ch *t.Channel, sampleRate int) {
	freq1 := ch.Voice.Carrier + ch.Voice.Resonance/2
	freq2 := ch.Voice.Carrier - ch.Voice.Resonance/2
	ch.Amplitude[0] = int(ch.Voice.Amplitude)
	ch.Amplitude[1] = int(ch.Voice.Amplitude)
	ch.Increment[0] = int(freq1 / float64(sampleRate) * t.SineTableSize * 65536)
	ch.Increment[1] = int(freq2 / float64(sampleRate) * t.SineTableSize * 65536)
}

// binauralGenerateSample generates a binaural beat sample
func binauralGenerateSample(ch *t.Channel, waveTables [4][]int) (int, int) {
	// Advance offset for each ear
	ch.Offset[0] += ch.Increment[0]
	ch.Offset[0] &= (t.SineTableSize << 16) - 1

	ch.Offset[1] += ch.Increment[1]
	ch.Offset[1] &= (t.SineTableSize << 16) - 1

	waveIdx := int(ch.Voice.Waveform)

	leftSample := ch.Amplitude[0] * waveTables[waveIdx][ch.Offset[0]>>16]
	rightSample := ch.Amplitude[1] * waveTables[waveIdx][ch.Offset[1]>>16]

	return leftSample, rightSample
}
