package shared

import (
	t "github.com/ruanklein/synapseq/internal/types"
)

// CountActiveChannels counts the number of active channels
func CountActiveChannels(chs []t.Channel) int {
	for i := len(chs) - 1; i >= 0; i-- {
		if chs[i].Track.Type != t.TrackOff {
			return i + 1
		}
	}
	return 1 // At least 1 channel always
}
