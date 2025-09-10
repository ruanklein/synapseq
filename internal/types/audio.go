/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

const (
	BufferSize         = 1024    // Buffer size for audio processing
	SineTableSize      = 16384   // Number of elements in sine-table (power of 2)
	WaveTableAmplitude = 0x7FFFF // Amplitude of wave in wave-table
	PhasePrecision     = 65536   // Phase precision (1/65536 of a cycle)
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
