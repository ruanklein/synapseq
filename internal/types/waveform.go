/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

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
