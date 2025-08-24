package sequence

import (
	"fmt"
)

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) ([]Preset, error) {
	file, err := LoadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	presets := make([]Preset, 0, MaxPresets)

	for file.NextLine() {
		ctx := NewLineContext(file.CurrentLine)

		// Skip empty lines
		if len(ctx.Tokens) == 0 {
			continue
		}

		// Skip comments
		if isCommentLine(ctx) {
			// Print the comment if it has more than 2 characters
			comment := parseCommentLine(ctx)
			if comment != "" {
				fmt.Printf("> %s\n", comment)
			}
			continue
		}

		// Preset definition
		if isValidPresetName(ctx.Line) {
			if len(presets) >= MaxPresets {
				return nil, fmt.Errorf("line %d: maximum number of presets reached", file.CurrentLineNumber)
			}

			preset, err := parsePresetLine(ctx)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			for _, p := range presets {
				if p.Name == preset.Name {
					return nil, fmt.Errorf("line %d: duplicate preset definition: %s", file.CurrentLineNumber, preset.Name)
				}
			}

			presets = append(presets, *preset)
			continue
		}

		// Voice line
		if isVoiceLine(ctx) {
			if len(presets) == 0 {
				return nil, fmt.Errorf("line %d: definition defined before any preset: %s", file.CurrentLineNumber, ctx.Line)
			}

			lastPreset := &presets[len(presets)-1]
			voiceIndex, err := lastPreset.AllocateVoice()
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			voice, err := parseVoiceLine(ctx)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			lastPreset.Voice[voiceIndex] = *voice
			continue
		}

		return nil, fmt.Errorf("line %d: invalid syntax: %s", file.CurrentLineNumber, ctx.Line)
	}

	// Check for empty presets
	for _, p := range presets {
		if p.AllVoicesAreOff() {
			return nil, fmt.Errorf("preset '%s' is empty", p.Name)
		}
	}

	return presets, nil
}
