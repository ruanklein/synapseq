package shared

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// Equal checks if two voice configurations are identical
func IsVoiceEqual(v1, v2 *t.Voice) bool {
	return v1.Type == v2.Type &&
		v1.Amplitude == v2.Amplitude &&
		v1.Carrier == v2.Carrier &&
		v1.Resonance == v2.Resonance &&
		v1.Waveform == v2.Waveform &&
		v1.Intensity == v2.Intensity
}
