/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

import "fmt"

// TransitionCurveK is the curve constant for logarithmic, exponential, and sigmoid transitions
const TransitionCurveK = 6.0

// TransitionType defines the type of slide for track transitions
type TransitionType int

const (
	TransitionSteady TransitionType = iota
	TransitionEaseOut
	TransitionEaseIn
	TransitionSmooth
)

// String returns the string representation of the TransitionType
func (s TransitionType) String() string {
	switch s {
	case TransitionSteady:
		return "steady"
	case TransitionEaseOut:
		return "ease-out"
	case TransitionEaseIn:
		return "ease-in"
	case TransitionSmooth:
		return "smooth"
	default:
		return "unknown"
	}
}

// Period represents a time period with track configurations
type Period struct {
	Time       int                     // Start time (end time is ->Next->Time)
	TrackStart [NumberOfChannels]Track // Start tracks for each channel
	TrackEnd   [NumberOfChannels]Track // End tracks for each channel
	Transition TransitionType          // Transition type
}

// TimeString returns the time of this period as a formatted string
func (p *Period) TimeString() string {
	hh := p.Time / 3600000
	mm := (p.Time % 3600000) / 60000
	ss := (p.Time % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
