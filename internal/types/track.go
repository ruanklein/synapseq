package types

import "fmt"

// TrackType represents the type of track/sound
type TrackType int

const (
	// Track is off
	TrackOff TrackType = iota
	// Track is silence
	TrackSilence
	// Track is a binaural beat
	TrackBinauralBeat
	// Track is a monaural beat
	TrackMonauralBeat
	// Track is an isochronic beat
	TrackIsochronicBeat
	// Track is white noise
	TrackWhiteNoise
	// Track is pink noise
	TrackPinkNoise
	// Track is brown noise
	TrackBrownNoise
	// Track is a spin white noise
	TrackSpinWhite
	// Track is a spin pink noise
	TrackSpinPink
	// Track is a spin brown noise
	TrackSpinBrown
	// Track is a background noise
	TrackBackground
)

// String returns the string representation of the TrackType
func (tr TrackType) String() string {
	switch tr {
	case TrackOff:
		return KeywordOff
	case TrackSilence:
		return KeywordSilence
	case TrackBinauralBeat:
		return KeywordBinaural
	case TrackMonauralBeat:
		return KeywordMonaural
	case TrackIsochronicBeat:
		return KeywordIsochronic
	case TrackWhiteNoise, TrackSpinWhite:
		return KeywordWhite
	case TrackPinkNoise, TrackSpinPink:
		return KeywordPink
	case TrackBrownNoise, TrackSpinBrown:
		return KeywordBrown
	case TrackBackground:
		return KeywordBackground
	default:
		return "unknown"
	}
}

// Track represents a track configuration
type Track struct {
	Type      TrackType     // Track type
	Amplitude AmplitudeType // Amplitude level (0-4096 for 0-100%)
	Carrier   float64       // Carrier frequency
	Resonance float64       // Resonance frequency
	Waveform  WaveformType  // Waveform shape
	Intensity IntensityType // Intensity (for background effects)
}

// Validate checks if the track configuration is valid
func (tr *Track) Validate() error {
	if tr.Amplitude < 0 || tr.Amplitude > 4096 {
		return fmt.Errorf("amplitude must be between 0 and 100. Received: %.2f", tr.Amplitude.ToPercent())
	}
	if tr.Carrier < 0 {
		return fmt.Errorf("carrier frequency must be positive. Received: %.2f", tr.Carrier)
	}
	if tr.Resonance < 0 {
		return fmt.Errorf("resonance frequency must be positive. Received: %.2f", tr.Resonance)
	}
	if tr.Intensity < 0 || tr.Intensity > 1.0 {
		return fmt.Errorf("intensity must be between 0 and 100. Received: %.2f", tr.Intensity.ToPercent())
	}
	return nil
}

// String returns the string representation of the Track configuration
func (tr *Track) String() string {
	switch tr.Type {
	case TrackOff, TrackSilence:
		return "- -"
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		return fmt.Sprintf("%s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		return fmt.Sprintf("%s %s %s %.2f", KeywordNoise, tr.Type.String(), KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackSpinWhite, TrackSpinPink, TrackSpinBrown:
		return fmt.Sprintf("%s %s %s %s %s %.2f %s %.2f %s %.2f", KeywordWaveform, tr.Waveform.String(), KeywordSpin, tr.Type.String(), KeywordWidth, tr.Carrier, KeywordRate, tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackBackground:
		return fmt.Sprintf("%s %s %.2f", KeywordBackground, KeywordAmplitude, tr.Amplitude.ToPercent())
	default:
		return ""
	}
}

// ShortString returns a compact string representation of the track configuration
func (tr *Track) ShortString() string {
	switch tr.Type {
	case TrackOff, TrackSilence:
		return " -"
	case TrackBinauralBeat, TrackMonauralBeat, TrackIsochronicBeat:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordTone, tr.Carrier, tr.Type.String(), tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackWhiteNoise, TrackPinkNoise, TrackBrownNoise:
		return fmt.Sprintf(" (%s:%.2f)", KeywordNoise, tr.Amplitude.ToPercent())
	case TrackSpinWhite, TrackSpinPink, TrackSpinBrown:
		return fmt.Sprintf(" (%s:%.2f %s:%.2f %s:%.2f)",
			KeywordSpin, tr.Carrier, KeywordRate, tr.Resonance, KeywordAmplitude, tr.Amplitude.ToPercent())
	case TrackBackground:
		return fmt.Sprintf(" (%s:%.2f)", KeywordAmplitude, tr.Amplitude.ToPercent())
	default:
		return " ???"
	}
}
