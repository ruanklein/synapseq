package sequence

import (
	"fmt"

	"github.com/ruanklein/synapseq/internal/parser"
	t "github.com/ruanklein/synapseq/internal/types"
)

// LoadSequence loads a sequence from a file
func LoadSequence(fileName string) ([]t.Period, *t.Option, error) {
	file, err := LoadFile(fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	presets := make([]t.Preset, 0, t.MaxPresets)
	var periods []t.Period

	// Initialize built-in presets
	silencePreset := t.Preset{Name: t.BuiltinSilence}
	silencePreset.InitVoices(t.VoiceSilence)
	presets = append(presets, silencePreset)

	// Initialize audio options
	options := &t.Option{
		SampleRate:     44100,
		Volume:         100,
		BackgroundPath: "",
		GainLevel:      t.GainLevelMedium,
	}

	// Parse each line in the file
	for file.NextLine() {
		ctx := parser.NewTextParser(file.CurrentLine)

		// Skip empty lines
		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		// Skip comments
		if ctx.HasComment() {
			comment := ctx.ParseComment()
			if comment != "" {
				fmt.Printf("> %s\n", comment)
			}
			continue
		}

		// Option line
		if ctx.HasOption() {
			if len(presets) > 1 {
				return nil, nil, fmt.Errorf("line %d: options must be defined before any preset", file.CurrentLineNumber)
			}

			if err = ctx.ParseOption(options); err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}
			continue
		}

		// Preset definition
		if ctx.HasPreset() {
			if len(presets) >= t.MaxPresets {
				return nil, nil, fmt.Errorf("line %d: maximum number of presets reached", file.CurrentLineNumber)
			}

			if len(periods) > 0 {
				return nil, nil, fmt.Errorf("line %d: preset definitions must be before any timeline definitions", file.CurrentLineNumber)
			}

			preset, err := ctx.ParsePreset()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			if preset.Name == t.BuiltinSilence {
				return nil, nil, fmt.Errorf("line %d: preset name %q is reserved", file.CurrentLineNumber, t.BuiltinSilence)
			}

			p := t.FindPreset(preset.Name, presets)
			if p != nil {
				return nil, nil, fmt.Errorf("line %d: duplicate preset definition: %s", file.CurrentLineNumber, preset.Name)
			}

			presets = append(presets, *preset)
			continue
		}

		// Voice line
		if ctx.HasVoice() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, nil, fmt.Errorf("line %d: voice defined before any preset: %s", file.CurrentLineNumber, ctx.Line.Raw)
			}

			if len(periods) > 0 {
				return nil, nil, fmt.Errorf("line %d: voice definitions must be before any timeline definitions", file.CurrentLineNumber)
			}

			lastPreset := &presets[len(presets)-1]
			voiceIndex, err := lastPreset.AllocateVoice()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			voice, err := ctx.ParseVoice()
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			lastPreset.Voice[voiceIndex] = *voice
			continue
		}

		// Timeline
		if ctx.HasTimeline() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, nil, fmt.Errorf("line %d: timeline defined before any preset: %s", file.CurrentLineNumber, ctx.Line.Raw)
			}

			period, err := ctx.ParseTimeline(&presets)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			if len(periods) == 0 && period.Time != 0 {
				return nil, nil, fmt.Errorf("line %d: first timeline must start at 00:00:00", file.CurrentLineNumber)
			}

			if len(periods) > 0 {
				lastPeriod := &periods[len(periods)-1]

				if lastPeriod.Time >= period.Time {
					return nil, nil, fmt.Errorf("line %d: timeline %s overlaps with previous timeline %s", file.CurrentLineNumber, period.TimeString(), lastPeriod.TimeString())
				}

				if err := t.AdjustPeriods(lastPeriod, period); err != nil {
					return nil, nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
				}
			}

			periods = append(periods, *period)
			continue
		}

		return nil, nil, fmt.Errorf("line %d: invalid syntax: %s", file.CurrentLineNumber, ctx.Line.Raw)
	}

	// Validate if has one preset (1 = silence preset)
	if len(presets) == 1 {
		return nil, nil, fmt.Errorf("no presets defined")
	}

	// Validate if all presets are defined
	for i := 1; i < len(presets); i++ {
		if presets[i].AllVoicesAreOff() {
			return nil, nil, fmt.Errorf("preset %q is empty", presets[i].Name)
		}
	}

	// Validate if has more than two Periods
	if len(periods) < 2 {
		return nil, nil, fmt.Errorf("at least two periods must be defined")
	}

	return periods, options, nil
}
