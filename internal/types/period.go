package types

import "fmt"

// Period represents a time period with voice configurations
type Period struct {
	Time       int                     // Start time (end time is ->Next->Time)
	VoiceStart [NumberOfChannels]Voice // Start voices for each channel
	VoiceEnd   [NumberOfChannels]Voice // End voices for each channel
}

// TimeString returns the time of this period as a formatted string
func (p *Period) TimeString() string {
	hh := p.Time / 3600000
	mm := (p.Time % 3600000) / 60000
	ss := (p.Time % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}
