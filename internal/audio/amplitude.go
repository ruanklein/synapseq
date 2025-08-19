package audio

// AmplitudePercentage represents a percentage value for amplitude
type AmplitudePercentage float64

// ToValue converts AmplitudePercentage to an integer value
func (ap AmplitudePercentage) ToValue() int {
	return int(ap * 40.96)
}

// FromValue converts an integer value to AmplitudePercentage
func (ap AmplitudePercentage) FromValue(value int) AmplitudePercentage {
	return AmplitudePercentage(float64(value) / 40.96)
}
