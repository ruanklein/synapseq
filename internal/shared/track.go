package shared

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// Equal checks if two track configurations are identical
func IsTrackEqual(tr1, tr2 *t.Track) bool {
	return tr1.Type == tr2.Type &&
		tr1.Amplitude == tr2.Amplitude &&
		tr1.Carrier == tr2.Carrier &&
		tr1.Resonance == tr2.Resonance &&
		tr1.Waveform == tr2.Waveform &&
		tr1.Intensity == tr2.Intensity
}
