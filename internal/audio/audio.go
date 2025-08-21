package audio

const (
	NumberOfChannels = 16    // Number of channels
	SineTableSize    = 16384 // Number of elements in sine-table (power of 2)

	amplitudeScale = 40.96 // Amplitude scale factor
)

var (
	BackgroundAmplitude *AmplitudeType // Background amplitude level to use with effects
)

type AmplitudeType float64 // Amplitude level (0-4096 for 0-100%)

// ToPercent converts a raw amplitude value to a float64 percentage
func (a AmplitudeType) ToPercent() float64 {
	return float64(a / amplitudeScale)
}

// AmplitudePercentToRaw converts a float64 value to a raw amplitude value
func AmplitudePercentToRaw(v float64) AmplitudeType {
	return AmplitudeType(v * amplitudeScale)
}
