package audio

import (
	"fmt"
)

type WaveformType int // WaveformType represents the waveform shape

// Waveform types
const (
	WaveformSine     WaveformType = iota // Sine
	WaveformSquare                       // Square
	WaveformTriangle                     // Triangle
	WaveformSawtooth                     // Sawtooth
)

// String returns the string representation of WaveformType:
// "sine", "square", "triangle", "sawtooth"
func (wt WaveformType) String() (string, error) {
	names := []string{"sine", "square", "triangle", "sawtooth"}

	if int(wt) >= len(names) {
		return "", fmt.Errorf("unknown waveform type: %d", wt)
	}

	return names[wt], nil
}
