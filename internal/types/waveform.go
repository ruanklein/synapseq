package types

// WaveformType represents the waveform shape
type WaveformType int

// Waveform types
const (
	WaveformSine     WaveformType = iota // Sine
	WaveformSquare                       // Square
	WaveformTriangle                     // Triangle
	WaveformSawtooth                     // Sawtooth
)

// String returns the string representation of WaveformType
func (wt WaveformType) String() string {
	switch wt {
	case WaveformSine:
		return KeywordSine
	case WaveformSquare:
		return KeywordSquare
	case WaveformTriangle:
		return KeywordTriangle
	case WaveformSawtooth:
		return KeywordSawtooth
	default:
		return ""
	}
}
