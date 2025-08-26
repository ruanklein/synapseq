// Audio processing types
package types

import (
	"fmt"
	"os"
)

const (
	NumberOfChannels   = 16      // Number of channels
	SineTableSize      = 16384   // Number of elements in sine-table (power of 2)
	WaveTableAmplitude = 0x7FFFF // Amplitude of wave in wave-table
)

// AudioOptions holds the configuration options for audio processing
type AudioOptions struct {
	SampleRate     int       // Sample rate (e.g., 44100)
	Volume         int       // Volume level (0-100 for 0-100%)
	BackgroundPath string    // Path to the background audio file
	GainLevel      GainLevel // Gain level (20, 16, 12, 6, 0) for audio processing
}

// Validate checks if the options are valid
func (o *AudioOptions) Validate() error {
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

type AudioFormat int // Audio format type

const (
	WavFormat AudioFormat = iota // WAV format
	RawFormat                    // RAW format
)

type GainLevel int // Gain level (-20db, -16db, -12db, -6db, 0db) for background audio

const (
	GainLevelVeryLow  GainLevel = 20 // -20db apply to background audio
	GainLevelLow      GainLevel = 16 // -16db apply to background audio
	GainLevelMedium   GainLevel = 12 // -12db apply to background audio
	GainLevelHigh     GainLevel = 6  // -6db apply to background audio
	GainLevelVeryHigh GainLevel = 0  // 0db apply to background audio
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
