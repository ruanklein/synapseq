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
			trackIndex, err := s.AllocateTrack(lastPreset)
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}

			track, err := ctx.ParseTrack()
			if err != nil {
				return nil, fmt.Errorf("preset file, line %d: %v", f.CurrentLineNumber, err)
			}

			if track.Type == t.TrackBackground {
				return nil, fmt.Errorf("preset file, line %d: background is not allowed in preset file", f.CurrentLineNumber)
			}

			lastPreset.Track[trackIndex] = *track
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
	}

	return presets, nil
}
