/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package shared

import (
	t "github.com/ruanklein/synapseq/v3/internal/types"
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
