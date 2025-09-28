/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"fmt"

	"github.com/ruanklein/synapseq/internal/parser"
	s "github.com/ruanklein/synapseq/internal/shared"
	t "github.com/ruanklein/synapseq/internal/types"
)

// LoadResult holds the result of loading a sequence
type LoadResult struct {
	Periods  []t.Period
	Options  *t.Option
	Comments []string
}

// LoadTextSequence loads a sequence from a text file
func LoadTextSequence(fileName string) (*LoadResult, error) {
	file, err := LoadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error loading sequence file: %v", err)
	}

	presets := make([]t.Preset, 0, t.MaxPresets)
	// For preset loading from file
	loadedPresets := false

	// Initialize built-in presets
	presets = append(presets, *t.NewBuiltinSilencePreset())

	// Initialize audio options
	options := &t.Option{
		SampleRate:     44100,
		Volume:         100,
		BackgroundPath: "",
		GainLevel:      t.GainLevelMedium,
	}

	var (
		periods  []t.Period
		comments []string
	)

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
				comments = append(comments, comment)
				// fmt.Fprintf(os.Stderr, "> %s\n", comment)
			}
			continue
		}

		// Option line
		if ctx.HasOption() {
			if len(presets) > 1 {
				return nil, fmt.Errorf("line %d: options must be defined before any preset", file.CurrentLineNumber)
			}

			if err = ctx.ParseOption(options); err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}
			// Validate options
			if err = options.Validate(); err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			// Load presets from file if specified in options and not already loaded
			if options.PresetPath != "" && !loadedPresets {
				fpresets, err := loadPresets(options.PresetPath)
				if err != nil {
					return nil, fmt.Errorf("%v", err)
				}
				presets = append(presets, fpresets...)
				loadedPresets = true
			}

			continue
		}

		// Preset definition
		if ctx.HasPreset() {
			if len(presets) >= t.MaxPresets {
				return nil, fmt.Errorf("line %d: maximum number of presets reached", file.CurrentLineNumber)
			}

			if len(periods) > 0 {
				return nil, fmt.Errorf("line %d: preset definitions must be before any timeline definitions", file.CurrentLineNumber)
			}

			preset, err := ctx.ParsePreset()
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			pName := preset.String()
			p := s.FindPreset(pName, presets)
			if p != nil {
				return nil, fmt.Errorf("line %d: duplicate preset definition: %s", file.CurrentLineNumber, pName)
			}

			presets = append(presets, *preset)
			continue
		}

		// Track line
		if ctx.HasTrack() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, fmt.Errorf("line %d: track defined before any preset: %s", file.CurrentLineNumber, ctx.Line.Raw)
			}

			if len(periods) > 0 {
				return nil, fmt.Errorf("line %d: track definitions must be before any timeline definitions", file.CurrentLineNumber)
			}

			lastPreset := &presets[len(presets)-1]
			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			if track.Type == t.TrackBackground && options.BackgroundPath == "" {
				return nil, fmt.Errorf("line %d: background track defined but no background audio file specified in options", file.CurrentLineNumber)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		// Timeline
		if ctx.HasTimeline() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, fmt.Errorf("line %d: timeline defined before any preset: %s", file.CurrentLineNumber, ctx.Line.Raw)
			}

			period, err := ctx.ParseTimeline(&presets)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
			}

			if len(periods) == 0 && period.Time != 0 {
				return nil, fmt.Errorf("line %d: first timeline must start at 00:00:00", file.CurrentLineNumber)
			}

			if len(periods) > 0 {
				lastPeriod := &periods[len(periods)-1]

				if lastPeriod.Time >= period.Time {
					return nil, fmt.Errorf("line %d: timeline %s overlaps with previous timeline %s", file.CurrentLineNumber, period.TimeString(), lastPeriod.TimeString())
				}

				if err := s.AdjustPeriods(lastPeriod, period); err != nil {
					return nil, fmt.Errorf("line %d: %v", file.CurrentLineNumber, err)
				}
			}

			periods = append(periods, *period)
			continue
		}

		return nil, fmt.Errorf("line %d: invalid syntax: %s", file.CurrentLineNumber, ctx.Line.Raw)
	}

	// Validate if has one preset (1 = silence preset)
	if len(presets) == 1 {
		return nil, fmt.Errorf("no presets defined")
	}

	if len(presets) > t.MaxPresets {
		return nil, fmt.Errorf("maximum number of presets exceeded: %d > %d", len(presets), t.MaxPresets)
	}

	// Validate each preset (skip silence preset)
	for i := 1; i < len(presets); i++ {
		p := &presets[i]
		if s.IsPresetEmpty(p) {
			return nil, fmt.Errorf("preset %q is empty", presets[i].String())
		}
		if n := s.NumBackgroundTracks(p); n > 1 {
			return nil, fmt.Errorf("preset %q has %d background tracks; only one background track is allowed per preset", presets[i].String(), n)
		}
	}

	// Validate if has more than two Periods
	if len(periods) < 2 {
		return nil, fmt.Errorf("at least two periods must be defined")
	}

	return &LoadResult{
		Periods:  periods,
		Options:  options,
		Comments: comments,
	}, nil
}
