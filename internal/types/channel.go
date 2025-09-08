package types

const (
	NumberOfChannels = 16 // Max number of channels
)

// Channel represents a channel state
type Channel struct {
	// Current track setting (updated from current period)
	Track Track
	// Track type
	Type TrackType
	// Current amplitude state
	Amplitude [2]int
	// Increment (for binaural tones, offset + increment into sine table * 65536)
	Increment [2]int
	// Offset into waveform table (for tones, offset + increment into sine table * 65536)
	Offset [2]int
}
