/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"fmt"
	"io"
	"strings"

	s "github.com/ruanklein/synapseq/internal/shared"
	t "github.com/ruanklein/synapseq/internal/types"
)

// StatusReporter handles terminal status output during rendering
type StatusReporter struct {
	// Output writer
	out io.Writer
	// To clear the previous line
	lastStatusWidth int
	// To detect period change
	lastPeriodIdx int
	// To control update frequency
	updateCounter int
}

// NewStatusReporter creates a new status reporter
func NewStatusReporter(out io.Writer) *StatusReporter {
	return &StatusReporter{
		out:           out,
		lastPeriodIdx: -1,
	}
}

// DisplayPeriodChange shows details of the period when it changes (like dispCurrPer)
func (sr *StatusReporter) DisplayPeriodChange(r *AudioRenderer, periodIdx int) {
	if periodIdx >= len(r.periods) || sr.out == nil {
		return
	}

	period := r.periods[periodIdx]
	var nextPeriod *t.Period
	if periodIdx+1 < len(r.periods) {
		nextPeriod = &r.periods[periodIdx+1]
	} else {
		// Last period - use the same as end
		nextPeriod = &period
	}

	// Clear previous line if necessary
	if sr.lastStatusWidth > 0 {
		fmt.Fprintf(sr.out, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		sr.lastStatusWidth = 0
	}

	// Line 1: Current period (start)
	line1 := fmt.Sprintf("- %s -> %s (%s)",
		period.TimeString(),
		nextPeriod.TimeString(),
		period.Transition.String())

	// Line 2: Start tracks (indented)
	line2 := ""

	for ch := range s.CountActiveChannels(r.channels[:]) {
		startTrack := period.TrackStart[ch]
		endTrack := period.TrackEnd[ch]

		// Start Track
		if startTrack.Type != t.TrackOff && startTrack.Type != t.TrackSilence {
			line2 += fmt.Sprintf("\n%s %s", strings.Repeat(" ", 6), startTrack.String())
		}

		// End Track (only if different)
		if !s.IsTrackEqual(&startTrack, &endTrack) {
			line2 += fmt.Sprintf("\n   ->  %s", endTrack.String())
		}
	}

	// Show the lines
	fmt.Fprintf(sr.out, "%s%s\n", line1, line2)
}

// DisplayStatus show the current status line
func (sr *StatusReporter) DisplayStatus(r *AudioRenderer, currentTimeMs int) {
	if sr.out == nil {
		return
	}

	// Format current time
	hh := currentTimeMs / 3600000
	mm := (currentTimeMs % 3600000) / 60000
	ss := (currentTimeMs % 60000) / 1000

	// Create status line
	status := fmt.Sprintf("  %02d:%02d:%02d", hh, mm, ss)

	// Add active tracks from each channel
	for ch := range s.CountActiveChannels(r.channels[:]) {
		channel := &r.channels[ch]
		status += channel.Track.ShortString()
	}

	// Clean previous line if necessary
	clearStr := ""
	if sr.lastStatusWidth > len(status) {
		clearStr = strings.Repeat(" ", sr.lastStatusWidth-len(status))
	}

	fmt.Fprintf(sr.out, "%s%s\r", status, clearStr)
	sr.lastStatusWidth = len(status)
}

// CheckPeriodChange checks if the period has changed and displays if necessary
func (sr *StatusReporter) CheckPeriodChange(r *AudioRenderer, periodIdx int) {
	if periodIdx != sr.lastPeriodIdx {
		sr.DisplayPeriodChange(r, periodIdx)
		sr.lastPeriodIdx = periodIdx
	}
}

// ShouldUpdateStatus checks if the status should be updated
func (sr *StatusReporter) ShouldUpdateStatus() bool {
	sr.updateCounter++
	// Update every ~44 buffers (~ 1 second at 44100Hz with buffer 1024)
	return sr.updateCounter%44 == 0
}

// FinalStatus clears the status line at the end
func (sr *StatusReporter) FinalStatus() {
	if sr.out == nil {
		return
	}
	if sr.lastStatusWidth > 0 {
		fmt.Fprintf(sr.out, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		fmt.Fprintf(sr.out, "\n")
	}
}
