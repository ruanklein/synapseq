package config

import "github.com/ruanklein/synapseq/internal/audio"

const (
	SampleRate   = 44100           // Sample rate
	OutputFormat = audio.WavFormat // Output format
	Volume       = 100             // Volume (0-100)
)
