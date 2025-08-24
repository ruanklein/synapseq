package audio

const (
	NumberOfChannels = 16    // Number of channels
	SineTableSize    = 16384 // Number of elements in sine-table (power of 2)
)

type AmplitudeType float64 // Amplitude level (0-4096 for 0-100%)

// ToPercent converts a raw amplitude value to a float64 percentage
func (a AmplitudeType) ToPercent() float64 {
	return float64(a / 40.96)
}

// AmplitudePercentToRaw converts a float64 value to a raw amplitude value
func AmplitudePercentToRaw(v float64) AmplitudeType {
	return AmplitudeType(v * 40.96)
}

type IntensityType float64 // Intensity level (0-1.0 for 0-100%)

// ToPercent converts a raw intensity value to a float64 percentage
func (i IntensityType) ToPercent() float64 {
	return float64(i * 100)
}

// IntensityPercentToRaw converts a float64 value to a raw intensity value
func IntensityPercentToRaw(v float64) IntensityType {
	return IntensityType(v / 100)
}
