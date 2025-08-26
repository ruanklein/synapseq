package types

var (
	CurrentPeriod *Period // Current period state
)

// Period represents a time period with voice configurations
type Period struct {
	Next       *Period                 // Next period in chain
	Prev       *Period                 // Previous period in chain
	Time       int                     // Start time (end time is ->Next->Time)
	VoiceStart [NumberOfChannels]Voice // Start voices for each channel
	VoiceEnd   [NumberOfChannels]Voice // End voices for each channel
	FadeIn     int                     // Fade-in mode
	FadeOut    int                     // Fade-out mode
}

// Duration returns the duration of this period in milliseconds
func (p *Period) Duration() int {
	if p.Next == nil {
		return 0 // Last period has no duration
	}
	return p.Next.Time - p.Time
}

// IsLast returns true if this is the last period in the chain
func (p *Period) IsLast() bool {
	return p.Next == nil
}

// IsFirst returns true if this is the first period in the chain
func (p *Period) IsFirst() bool {
	return p.Prev == nil
}
