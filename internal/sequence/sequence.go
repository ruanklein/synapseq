package sequence

import (
	"fmt"
	"os"
)

type GainLevel int // Gain level (-20db, -16db, -12db, -6db, 0db) for background audio

const (
	gainVeryLow  GainLevel = 20 // -20db apply to background audio
	gainLow      GainLevel = 16 // -16db apply to background audio
	gainMedium   GainLevel = 12 // -12db apply to background audio
	gainHigh     GainLevel = 6  // -6db apply to background audio
	gainVeryHigh GainLevel = 0  // 0db apply to background audio
)

type SequenceOptions struct {
	SampleRate     int       // Sample rate (e.g., 44100)
	Volume         int       // Volume level (0-100 for 0-100%)
	BackgroundPath string    // Path to the background audio file
	GainLevel      GainLevel // Gain level (20, 16, 12, 6, 0)
}

// Validate checks if the sequence options are valid
func (s *SequenceOptions) Validate() error {
	if s.SampleRate <= 0 {
		return fmt.Errorf("invalid sample rate: %d", s.SampleRate)
	}
	if s.Volume < 0 || s.Volume > 100 {
		return fmt.Errorf("invalid volume: %d", s.Volume)
	}
	if s.BackgroundPath != "" {
		if _, err := os.Stat(s.BackgroundPath); os.IsNotExist(err) {
			return fmt.Errorf("background path does not exist: %s", s.BackgroundPath)
		}
	}
	return nil
}

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) ([]Preset, *SequenceOptions, error) {
	file, err := LoadFile(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	presets := make([]Preset, 0, MaxPresets)

	// Initialize built-in presets
	silencePreset := Preset{Name: builtinSilence}
	silencePreset.InitVoices()
	presets = append(presets, silencePreset)

	// Initialize sequence options
	options := &SequenceOptions{
		SampleRate:     44100,
		Volume:         100,
		BackgroundPath: "",
		GainLevel:      gainMedium,
	}

	for file.NextLine() {
		ctx := NewLineContext(file.CurrentLine)

		// Skip empty lines
		if len(ctx.Tokens) == 0 {
			continue
		}

		// Skip comments
		if isCommentLine(ctx) {
			comment := parseCommentLine(ctx)
			if comment != "" {
				fmt.Printf("> %s\n", comment)
			}
			continue
		}

		// Option line
		if isOptionLine(ctx) {
			if len(presets) > 1 {
				return nil, nil, fmt.Errorf("line %d: options must be defined before any preset", file.CurrentLineNumber)
			}

			if err := parseOptionLine(ctx, options); err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}
			continue
		}

		// Preset definition
		if isValidPresetName(ctx.Line) {
			if len(presets) >= MaxPresets {
				return nil, nil, fmt.Errorf("line %d: maximum number of presets reached", file.CurrentLineNumber)
			}

			preset, err := parsePresetLine(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			for _, p := range presets {
				if p.Name == preset.Name {
					return nil, nil, fmt.Errorf("line %d: duplicate preset definition: %s", file.CurrentLineNumber, preset.Name)
				}
			}

			presets = append(presets, *preset)
			continue
		}

		// Voice line
		if isVoiceLine(ctx) {
			if len(presets) == 1 { // 1 = silence preset
				return nil, nil, fmt.Errorf("line %d: definition defined before any preset: %s", file.CurrentLineNumber, ctx.Line)
			}

			lastPreset := &presets[len(presets)-1]
			voiceIndex, err := lastPreset.AllocateVoice()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			voice, err := parseVoiceLine(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			lastPreset.Voice[voiceIndex] = *voice
			continue
		}

		return nil, nil, fmt.Errorf("line %d: invalid syntax: %s", file.CurrentLineNumber, ctx.Line)
	}

	// Check for empty presets
	for i := 1; i < len(presets); i++ {
		if presets[i].AllVoicesAreOff() {
			return nil, nil, fmt.Errorf("preset '%s' is empty", presets[i].Name)
		}
	}

	// Validate options
	if err := options.Validate(); err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	return presets, options, nil
}
