package sequence

import (
	"fmt"

	"github.com/ruanklein/synapseq/internal/parser"
	t "github.com/ruanklein/synapseq/internal/types"
)

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) ([]t.Preset, *t.AudioOptions, error) {
	file, err := LoadFile(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	presets := make([]t.Preset, 0, t.MaxPresets)

	// Initialize built-in presets
	silencePreset := t.Preset{Name: t.BuiltinSilence}
	silencePreset.InitVoices()
	presets = append(presets, silencePreset)

	// Initialize audio options
	options := &t.AudioOptions{
		SampleRate:     44100,
		Volume:         100,
		BackgroundPath: "",
		GainLevel:      t.GainLevelMedium,
	}

	for file.NextLine() {
		ctx := parser.NewParserContext(file.CurrentLine)

		// Skip empty lines
		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		// Skip comments
		if ctx.IsCommentLine() {
			comment := ctx.ParseCommentLine()
			if comment != "" {
				fmt.Printf("> %s\n", comment)
			}
			continue
		}

		// Option line
		if ctx.IsOptionLine() {
			if len(presets) > 1 {
				return nil, nil, fmt.Errorf("line %d: options must be defined before any preset", file.CurrentLineNumber)
			}

			if err := ctx.ParseOptionLine(options); err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}
			continue
		}

		// Preset definition
		if ctx.IsPresetLine() {
			if len(presets) >= t.MaxPresets {
				return nil, nil, fmt.Errorf("line %d: maximum number of presets reached", file.CurrentLineNumber)
			}

			preset, err := ctx.ParsePresetLine()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			if preset.Name == t.BuiltinSilence {
				return nil, nil, fmt.Errorf("line %d: preset name %q is reserved", file.CurrentLineNumber, t.BuiltinSilence)
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
		if ctx.IsVoiceLine() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, nil, fmt.Errorf("line %d: definition defined before any preset: %s", file.CurrentLineNumber, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			voiceIndex, err := lastPreset.AllocateVoice()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			voice, err := ctx.ParseVoiceLine()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			lastPreset.Voice[voiceIndex] = *voice
			continue
		}

		// Timeline (in development)
		if ctx.IsTimeline() {
			timeline, err := ctx.ParseTimeline()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}
			fmt.Print(timeline)
			continue
		}

		return nil, nil, fmt.Errorf("line %d: invalid syntax: %s", file.CurrentLineNumber, ctx.Line.Raw)
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
