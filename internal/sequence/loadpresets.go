/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package sequence

import (
	"fmt"

	"github.com/ruanklein/synapseq/v3/internal/parser"
	s "github.com/ruanklein/synapseq/v3/internal/shared"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// loadPresets loads presets from a given file path
func loadPresets(filename string) ([]t.Preset, error) {
	f, err := LoadFile(filename)
	if err != nil {
		return nil, err
	}

	presets := make([]t.Preset, 0, t.MaxPresets)
	for f.NextLine() {
		ctx := parser.NewTextParser(f.CurrentLine)

		// Skip empty lines
		if len(ctx.Line.Tokens) == 0 {
			continue
		}

		// Skip comments
		if ctx.HasComment() {
			continue
		}

		// Parse preset lines
		if ctx.HasPreset() {
			preset, err := ctx.ParsePreset(&presets)
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}
			presets = append(presets, *preset)
			continue
		}

		// Track line
		if ctx.HasTrack() {
			if len(presets) == 0 {
				return nil, fmt.Errorf("preset file, line %d: track defined before any preset: %s", f.CurrentLineNumber, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.From != nil {
				return nil, fmt.Errorf("preset file, line %d: preset %q inherits from another and cannot define new tracks", f.CurrentLineNumber, lastPreset.String())
			}

			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}

			lastPreset.Track[trackIndex] = *track
			continue
		}

		// Track override line
		if ctx.HasTrackOverride() {
			if len(presets) == 1 { // 1 = silence preset
				return nil, fmt.Errorf("preset file, line %d: track override defined before any preset: %s", f.CurrentLineNumber, ctx.Line.Raw)
			}

			lastPreset := &presets[len(presets)-1]
			if lastPreset.IsTemplate {
				return nil, fmt.Errorf("preset file, line %d: cannot override tracks on template preset %q", f.CurrentLineNumber, lastPreset.String())
			}
			if lastPreset.From == nil {
				return nil, fmt.Errorf("preset file, line %d: cannot override tracks on preset %q which does not have a 'from' source", f.CurrentLineNumber, lastPreset.String())
			}

			if err := ctx.ParseTrackOverride(lastPreset); err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}

			continue
		}

		return nil, fmt.Errorf("preset file, line %d: unexpected content: %s", f.CurrentLineNumber, f.CurrentLine)
	}

	// Validate if has one preset
	if len(presets) == 0 {
		return nil, fmt.Errorf("preset file: no presets defined")
	}

	// Validate each preset (skip silence preset)
	for _, p := range presets {
		if s.IsPresetEmpty(&p) {
			return nil, fmt.Errorf("preset file: preset %q is empty", p.String())
		}
		if n := s.NumBackgroundTracks(&p); n > 1 {
			return nil, fmt.Errorf("preset file: preset %q has %d background tracks; only one background track is allowed per preset", p.String(), n)
		}
	}

	return presets, nil
}
