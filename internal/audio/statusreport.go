package audio

import (
	"fmt"
	"os"
	"strings"

	s "github.com/ruanklein/synapseq/internal/shared"
	t "github.com/ruanklein/synapseq/internal/types"
)

// StatusReporter handles terminal status output during rendering
type StatusReporter struct {
	// If true, suppresses all output
	quiet bool
	// To clear the previous line
	lastStatusWidth int
	// To detect period change
	lastPeriodIdx int
	// To control update frequency
	updateCounter int
}

// NewStatusReporter creates a new status reporter
func NewStatusReporter(quiet bool) *StatusReporter {
	return &StatusReporter{
		quiet:         quiet,
		lastPeriodIdx: -1,
	}
}

// DisplayPeriodChange shows details of the period when it changes (like dispCurrPer)
func (sr *StatusReporter) DisplayPeriodChange(r *AudioRenderer) {
	if sr.quiet {
		return
	}

	if r.periodIdx >= len(r.periods) {
		return
	}

	period := r.periods[r.periodIdx]
	var nextPeriod *t.Period
	if r.periodIdx+1 < len(r.periods) {
		nextPeriod = &r.periods[r.periodIdx+1]
	} else {
		// Last period - use the same as end
		nextPeriod = &period
	}

	// Clear previous line if necessary
	if sr.lastStatusWidth > 0 {
		fmt.Fprintf(os.Stderr, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		sr.lastStatusWidth = 0
	}

	// Line 1: Current period (start)
	line1 := fmt.Sprintf("- %s ", period.TimeString())

	// Line 2: Next period (end)
	line2 := fmt.Sprintf("  %s ", nextPeriod.TimeString())

	for ch := range s.CountActiveChannels(r.channels[:]) {
		startTrack := period.TrackStart[ch]
		endTrack := period.TrackEnd[ch]

		// Start Track
		startStr := ""
		if startTrack.Type != t.TrackOff && startTrack.Type != t.TrackSilence {
			startStr = fmt.Sprintf("\n%s %s", strings.Repeat(" ", 6), startTrack.String())
		}

		// End Track
		endStr := "\n       --"
		if !s.IsTrackEqual(&startTrack, &endTrack) {
			endStr = "\n       -"
			if endTrack.Type != t.TrackOff && endTrack.Type != t.TrackSilence {
				endStr = fmt.Sprintf("\n%s %s", strings.Repeat(" ", 6), endTrack.String())
			}
		}

		line1 += startStr
		line2 += endStr
	}

	// Show the lines
	fmt.Fprintf(os.Stderr, "%s\n%s\n", line1, line2)
}

// DisplayStatus show the current status line
func (sr *StatusReporter) DisplayStatus(r *AudioRenderer, currentTimeMs int) {
	if sr.quiet {
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

	fmt.Fprintf(os.Stderr, "%s%s\r", status, clearStr)
	sr.lastStatusWidth = len(status)
}

// CheckPeriodChange checks if the period has changed and displays if necessary
func (sr *StatusReporter) CheckPeriodChange(r *AudioRenderer) {
	if r.periodIdx != sr.lastPeriodIdx {
		sr.DisplayPeriodChange(r)
		sr.lastPeriodIdx = r.periodIdx
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
	if !sr.quiet && sr.lastStatusWidth > 0 {
		fmt.Fprintf(os.Stderr, "%s\r", strings.Repeat(" ", sr.lastStatusWidth))
		fmt.Fprintf(os.Stderr, "\n")
	}
}
