/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	t "github.com/ruanklein/synapseq/internal/types"
)

// getFullPath expands ~ to the user's home directory and returns the absolute path
func getFullPath(path string) (string, error) {
	var fullPath string
	if path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
		fullPath = strings.Replace(path, "~", homeDir, 1)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		fullPath = filepath.Join(cwd, path)
	}

	return filepath.Clean(fullPath), nil
}

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
		isRemote := strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://")

		if content == "-" {
			return fmt.Errorf("stdin (-) is not supported for background or preset list")
		}

		fullPath := content
		if !isRemote {
			var err error
			fullPath, err = getFullPath(content)
			if err != nil {
				return fmt.Errorf("path: %v", err)
			}
		}

		if option == t.KeywordBackground {
			options.BackgroundPath = fullPath
		} else {
			options.PresetPath = fullPath
		}
	case t.KeywordOptionGainLevel:
		gainLevel, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected gain level: %s", ln)
		}

		switch gainLevel {
		case t.KeywordOptionGainLevelVeryLow:
			options.GainLevel = t.GainLevelVeryLow
		case t.KeywordOptionGainLevelLow:
			options.GainLevel = t.GainLevelLow
		case t.KeywordOptionGainLevelMedium:
			options.GainLevel = t.GainLevelMedium
		case t.KeywordOptionGainLevelHigh:
			options.GainLevel = t.GainLevelHigh
		case t.KeywordOptionGainLevelVeryHigh:
			options.GainLevel = t.GainLevelVeryHigh
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
