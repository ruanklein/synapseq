package audio

type AudioFormat int

// Supported audio formats
const (
	WavFormat AudioFormat = iota
	RawFormat
)
