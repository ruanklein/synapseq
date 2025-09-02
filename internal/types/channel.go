package types

// Channel represents a channel state
type Channel struct {
	Voice     Voice     // Current voice setting (updated from current period)
	Type      VoiceType // Voice type
	Amplitude [2]int    // Current amplitude state.
	Increment [2]int    // Increment (for binaural tones, offset + increment into sine table * 65536)
	Offset    [2]int    // Offset 1
}

// CountActiveChannels counts the number of active channels
func CountActiveChannels(chs []Channel) int {
	for i := len(chs) - 1; i >= 0; i-- {
		if chs[i].Voice.Type != VoiceOff {
			return i + 1
		}
	}
	return 1 // At least 1 channel always
}
