/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

// FormatOptions holds the options for the sequence format
type FormatOptions struct {
	Samplerate int `json:"samplerate" xml:"samplerate"`
	Volume     int `json:"volume" xml:"volume"`
}

// FormatElement represents a single element in the sequence format
type FormatElement struct {
	Tones  []FormatToneElement  `json:"tones,omitempty" xml:"tones,omitempty"`
	Noises []FormatNoiseElement `json:"noises,omitempty" xml:"noises,omitempty"`
}

// FormatToneElement represents a tone element in the sequence format
type FormatToneElement struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty"`
	Carrier   float64 `json:"carrier,omitempty" xml:"carrier,attr,omitempty"`
	Resonance float64 `json:"resonance,omitempty" xml:"resonance,attr,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty"`
	Waveform  string  `json:"waveform,omitempty" xml:"waveform,attr,omitempty"`
}

// FormatNoiseElement represents a noise element in the sequence format
type FormatNoiseElement struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty"`
}

// FormatSequenceEntry represents a single entry in the sequence format
type FormatSequenceEntry struct {
	Time     int           `json:"time" xml:"time,attr"`
	Elements FormatElement `json:"elements" xml:"elements>element"`
}

// SynapSeqInput represents the overall structure of a SynapSeq sequence file
type SynapSeqInput struct {
	Description []string              `json:"description" xml:"description>line"`
	Options     FormatOptions         `json:"options" xml:"options"`
	Sequence    []FormatSequenceEntry `json:"sequence" xml:"sequence>entry"`
}
