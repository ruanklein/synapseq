/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

// FormatInfo holds metadata about the sequence format
type FormatInfo struct {
	Name        string `json:"name" xml:"name"`
	Description string `json:"description" xml:"description"`
	Version     string `json:"version" xml:"version"`
}

// FormatOptions holds the options for the sequence format
type FormatOptions struct {
	Samplerate     int    `json:"samplerate" xml:"samplerate"`
	Volume         int    `json:"volume" xml:"volume"`
	BackgroundPath string `json:"background,omitempty" xml:"background,omitempty"`
	PresetPath     string `json:"presetpath,omitempty" xml:"presetpath,omitempty"`
	GainLevel      string `json:"gainlevel,omitempty" xml:"gainlevel,omitempty"`
}

// FormatElement represents a single element in the sequence format
type FormatElement struct {
	Kind      string  `json:"kind" xml:"kind,attr"`
	Mode      string  `json:"mode,omitempty" xml:"mode,attr,omitempty"`
	Carrier   float64 `json:"carrier,omitempty" xml:"carrier,attr,omitempty"`
	Resonance float64 `json:"resonance,omitempty" xml:"resonance,attr,omitempty"`
	Amplitude float64 `json:"amplitude,omitempty" xml:"amplitude,attr,omitempty"`
	Waveform  string  `json:"waveform,omitempty" xml:"waveform,attr,omitempty"`
}

// FormatSequenceEntry represents a single entry in the sequence format
type FormatSequenceEntry struct {
	Time  int             `json:"time" xml:"time,attr"`
	Track []FormatElement `json:"track" xml:"track>element"`
}

// SynapSeqInput represents the overall structure of a SynapSeq sequence file
type SynapSeqInput struct {
	Info     FormatInfo            `json:"info" xml:"info"`
	Options  FormatOptions         `json:"options" xml:"options"`
	Sequence []FormatSequenceEntry `json:"sequence" xml:"sequence>entry"`
}
