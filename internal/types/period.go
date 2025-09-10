/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

import "fmt"

// Period represents a time period with track configurations
type Period struct {
	Time       int                     // Start time (end time is ->Next->Time)
	TrackStart [NumberOfChannels]Track // Start tracks for each channel
	TrackEnd   [NumberOfChannels]Track // End tracks for each channel
}

// TimeString returns the time of this period as a formatted string
func (p *Period) TimeString() string {
	hh := p.Time / 3600000
	mm := (p.Time % 3600000) / 60000
	ss := (p.Time % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
