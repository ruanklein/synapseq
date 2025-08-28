package types

import (
	"fmt"
	"os"
)

// Option represents configuration options for audio processing.
type Option struct {
	SampleRate     int       // Sample rate (e.g., 44100)
	Volume         int       // Volume level (0-100 for 0-100%)
	BackgroundPath string    // Path to the background audio file
	GainLevel      GainLevel // Gain level (20, 16, 12, 6, 0) for audio processing
}

// Validate checks if the options are valid
func (o *Option) Validate() error {
	if o.SampleRate <= 0 {
		return fmt.Errorf("invalid sample rate: %d", o.SampleRate)
	}
	if o.Volume < 0 || o.Volume > 100 {
		return fmt.Errorf("invalid volume: %d", o.Volume)
	}
	if o.BackgroundPath != "" {
		if _, err := os.Stat(o.BackgroundPath); os.IsNotExist(err) {
			return fmt.Errorf("background path does not exist: %s", o.BackgroundPath)
		}
	}
	return nil
}

// String returns a string representation of the options.
func (o *Option) String() string {
	var background string

	if o.BackgroundPath == "" {
		background = "none"
	} else {
		background = o.BackgroundPath
	}

	return fmt.Sprintf("samplerate %d volume %d background %s gainlevel %ddb",
		o.SampleRate, o.Volume, background, o.GainLevel)
}
