package generator

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// monauralUpdateChannel updates channel state for monaural beats
func monauralUpdateChannel(ch *t.Channel, sampleRate int) {
	freqHigh := ch.Voice.Carrier + ch.Voice.Resonance/2
	freqLow := ch.Voice.Carrier - ch.Voice.Resonance/2
	ch.Amplitude[0] = int(ch.Voice.Amplitude)
	ch.Increment[0] = int(freqHigh / float64(sampleRate) * t.SineTableSize * 65536)
	ch.Increment[1] = int(freqLow / float64(sampleRate) * t.SineTableSize * 65536)
}

// monauralGenerateSample generates a monaural beat sample
func monauralGenerateSample(ch *t.Channel, waveTables [4][]int) (int, int) {
	// Advance phases for both frequencies
	ch.Offset[0] += ch.Increment[0] // high freq
	ch.Offset[0] &= (t.SineTableSize << 16) - 1

	ch.Offset[1] += ch.Increment[1] // low freq
	ch.Offset[1] &= (t.SineTableSize << 16) - 1

	waveIdx := int(ch.Voice.Waveform)

	freqHigh := waveTables[waveIdx][ch.Offset[0]>>16]
	freqLow := waveTables[waveIdx][ch.Offset[1]>>16]

	// Monaural: sum frequencies with reduced amplitude
	halfAmp := ch.Amplitude[0] / 2
	mixedSample := halfAmp * (freqHigh + freqLow)

	return mixedSample, mixedSample // Same content both ears
}
