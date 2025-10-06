/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

// FormatOptions holds the options for the sequence format
type FormatOptions struct {
	Samplerate int `json:"samplerate" xml:"samplerate" yaml:"samplerate"`
	Volume     int `json:"volume" xml:"volume" yaml:"volume"`
}

// FormatTrack represents a single element in the sequence format
type FormatTrack struct {
	Tones  []FormatToneTrack  `json:"tones,omitempty" xml:"tone,omitempty" yaml:"tones"`
	Noises []FormatNoiseTrack `json:"noises,omitempty" xml:"noise,omitempty" yaml:"noises"`
}

// FormatToneTrack represents a tone element in the sequence format
type FormatToneTrack struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty" yaml:"mode"`
	Carrier   float64 `json:"carrier,omitempty" xml:"carrier,attr,omitempty" yaml:"carrier"`
	Resonance float64 `json:"resonance,omitempty" xml:"resonance,attr,omitempty" yaml:"resonance"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty" yaml:"amplitude"`
	Waveform  string  `json:"waveform,omitempty" xml:"waveform,attr,omitempty" yaml:"waveform"`
}

// FormatNoiseTrack represents a noise element in the sequence format
type FormatNoiseTrack struct {
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty" yaml:"mode"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty" yaml:"amplitude"`
}

// FormatSequenceEntry represents a single entry in the sequence format
type FormatSequenceEntry struct {
	Time  int         `json:"time" xml:"time,attr" yaml:"time"`
	Track FormatTrack `json:"track" xml:"track" yaml:"track"`
}

// SynapSeqInput represents the overall structure of a SynapSeq sequence file
type SynapSeqInput struct {
	Description []string              `json:"description" xml:"description>line" yaml:"description"`
	Options     FormatOptions         `json:"options" xml:"options" yaml:"options"`
	Sequence    []FormatSequenceEntry `json:"sequence" xml:"sequence>entry" yaml:"sequence"`
}
