/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// HasTrackOverride checks if the current line is a track override definition
func (ctx *TextParser) HasTrackOverride() bool {
	ln := ctx.Line.Raw
	if len(ln) < 3 {
		return false
	}

	if ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' ' {
		tok, ok := ctx.Line.Peek()
		if !ok || tok != t.KeywordTrack {
			return false
		}
		return true
	}

	return false
}

// ParseTrackOverride applies track overrides to the given preset
func (ctx *TextParser) ParseTrackOverride(preset *t.Preset) error {
	if preset == nil || preset.From == nil {
		return fmt.Errorf("cannot override tracks on a preset without a 'from' source")
	}

	ln := ctx.Line.Raw
	_, ok := ctx.Line.NextToken()
	if !ok {
		return fmt.Errorf("expected 'track' keyword, got EOF: %s", ln)
	}

	trackIdx, err := ctx.Line.NextIntStrict()
	if err != nil {
		return fmt.Errorf("expected track index after 'track': %s", ln)
	}

	if trackIdx <= 0 || trackIdx >= t.NumberOfChannels {
		return fmt.Errorf("track index out of range (1-%d): %d", t.NumberOfChannels-1, trackIdx)
	}

	idx := trackIdx - 1 // Convert to 0-based index
	from := preset.From

	if from.Track[idx].Type == t.TrackOff {
		return fmt.Errorf("cannot override track %d which is off in the template preset %q", trackIdx, from.String())
	}

	kind, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordFrequency, t.KeywordAmplitude, t.KeywordIntensity)
	if err != nil {
		return fmt.Errorf("expected one of %q, %q, %q, %q: %s", t.KeywordTone, t.KeywordFrequency, t.KeywordAmplitude, t.KeywordIntensity, ln)
	}

	switch kind {
	case t.KeywordTone:
		carrier, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("carrier: %w", err)
		}

		preset.Track[idx].Carrier = carrier
	case t.KeywordFrequency:
		frequency, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("frequency: %w", err)
		}

		preset.Track[idx].Resonance = frequency
	case t.KeywordAmplitude:
		amplitude, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("amplitude: %w", err)
		}

		preset.Track[idx].Amplitude = t.AmplitudePercentToRaw(amplitude)
	case t.KeywordIntensity:
		intensity, err := ctx.Line.NextFloat64Strict()
		if err != nil {
			return fmt.Errorf("intensity: %w", err)
		}

		preset.Track[idx].Effect.Intensity = t.IntensityPercentToRaw(intensity)
	default:
		return fmt.Errorf("unexpected keyword: %s", kind)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return fmt.Errorf("unexpected token after track override definition: %q", unknown)
	}

	// Validate the updated track
	if err := preset.Track[idx].Validate(); err != nil {
		return fmt.Errorf("invalid track %d after override: %w", trackIdx, err)
	}

	return nil
}
