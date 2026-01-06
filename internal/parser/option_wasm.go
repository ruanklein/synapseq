//go:build wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

package parser

import (
	"fmt"
	"strings"

	s "github.com/synapseq-foundation/synapseq/v3/internal/shared"
	t "github.com/synapseq-foundation/synapseq/v3/internal/types"
)

// HasOption checks if the first element is an option
func (ctx *TextParser) HasOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and applies the option from the elements
func (ctx *TextParser) ParseOption(options *t.SequenceOptions) error {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return fmt.Errorf("expected option, got EOF: %s", ln)
	}

	if string(tok[0]) != t.KeywordOption {
		return fmt.Errorf("expected option. Received: %s", tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return fmt.Errorf("expected option name: %s", ln)
	}

	switch option {
	case t.KeywordOptionSampleRate:
		sampleRate, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("samplerate: %v", err)
		}
		options.SampleRate = sampleRate
	case t.KeywordOptionVolume:
		volume, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("volume: %v", err)
		}
		options.Volume = volume
	case t.KeywordOptionBackground, t.KeywordOptionPresetList:
		_, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected path: %s", ln)
		}

		content := strings.Join(ctx.Line.Tokens[1:], " ")
		if !s.IsRemoteFile(content) {
			return fmt.Errorf("file paths are not supported in WASM for background or preset list: %s", content)
		}

		if option == t.KeywordBackground {
			options.BackgroundPath = content
		} else {
			options.PresetList = append(options.PresetList, content)
		}
	case t.KeywordOptionGainLevel:
		gainLevel, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected gain level: %s", ln)
		}

		switch gainLevel {
		case t.KeywordOff:
			options.GainLevel = t.GainLevelOff
		case t.KeywordOptionGainLevelHigh:
			options.GainLevel = t.GainLevelHigh
		case t.KeywordOptionGainLevelMedium:
			options.GainLevel = t.GainLevelMedium
		case t.KeywordOptionGainLevelLow:
			options.GainLevel = t.GainLevelLow
		default:
			return fmt.Errorf("invalid gain level: %q", gainLevel)
		}
	default:
		return fmt.Errorf("invalid option: %q", option)
	}

	// If the option is not background, ensure no extra tokens are present
	if option != t.KeywordOptionBackground {
		unknown, ok := ctx.Line.Peek()
		if ok {
			return fmt.Errorf("unexpected token after option definition: %q", unknown)
		}
	}

	return nil
}
