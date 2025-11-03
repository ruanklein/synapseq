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

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// getFullPath resolves the full path of a given file path
func getFullPath(path, basePath string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		expanded := filepath.Join(homeDir, strings.TrimPrefix(path, "~"))
		return filepath.Clean(expanded), nil
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	fullPath := filepath.Join(basePath, path)
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
func (ctx *TextParser) ParseOption(options *t.SequenceOptions, filePath string) error {
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
			fullPath, err = getFullPath(content, filePath)
			if err != nil {
				return fmt.Errorf("path: %v", err)
			}
		}

		if option == t.KeywordBackground {
			options.BackgroundPath = fullPath
		} else {
			options.PresetList = append(options.PresetList, fullPath)
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
