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

// FormatTrack represents a single element in the sequence format
type FormatTrack struct {
	Tones  []FormatToneTrack  `json:"tones,omitempty" xml:"tone,omitempty"`
	Noises []FormatNoiseTrack `json:"noises,omitempty" xml:"noise,omitempty"`
}

// FormatToneTrack represents a tone element in the sequence format
type FormatToneTrack struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty"`
	Carrier   float64 `json:"carrier,omitempty" xml:"carrier,attr,omitempty"`
	Resonance float64 `json:"resonance,omitempty" xml:"resonance,attr,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty"`
	Waveform  string  `json:"waveform,omitempty" xml:"waveform,attr,omitempty"`
}

// FormatNoiseTrack represents a noise element in the sequence format
type FormatNoiseTrack struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty"`
}

// FormatSequenceEntry represents a single entry in the sequence format
type FormatSequenceEntry struct {
	Time  int         `json:"time" xml:"time,attr"`
	Track FormatTrack `json:"track" xml:"track"`
}

// SynapSeqInput represents the overall structure of a SynapSeq sequence file
type SynapSeqInput struct {
	Description []string              `json:"description" xml:"description>line"`
	Options     FormatOptions         `json:"options" xml:"options"`
	Sequence    []FormatSequenceEntry `json:"sequence" xml:"sequence>entry"`
}
