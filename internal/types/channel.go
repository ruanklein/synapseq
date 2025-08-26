package types

var (
	CurrentChannel [NumberOfChannels]Channel // Current channel state
)

// Channel represents a channel state
type Channel struct {
	Voice     Voice     // Current voice setting (updated from current period)
	Type      VoiceType // Voice type
	Amplitude [2]int    // Current amplitude state.
	Increment [2]int    // Increment (for binaural tones, offset + increment into sine table * 65536)
	Offset    [2]int    // Offset 1
}
