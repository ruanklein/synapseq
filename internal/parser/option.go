package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// IsOptionLine checks if the first element is an option
func (ctx *ParserContext) IsOptionLine() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == keywordOption
}

// ParseOption extracts and applies the option from the elements
func (ctx *ParserContext) ParseOptionLine(options *t.AudioOptions) error {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return fmt.Errorf("expected option, got EOF: %s", ln)
	}

	if string(tok[0]) != keywordOption {
		return fmt.Errorf("expected option. Received: %s", tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return fmt.Errorf("expected option name: %s", ln)
	}

	switch option {
	case keywordOptionSampleRate:
		sampleRate, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("samplerate: %v", err)
		}
		options.SampleRate = sampleRate
	case keywordOptionVolume:
		volume, err := ctx.Line.NextIntStrict()
		if err != nil {
			return fmt.Errorf("volume: %v", err)
		}
		options.Volume = volume
	case keywordOptionBackground:
		backgroundPath, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected background path: %s", ln)
		}
		options.BackgroundPath = backgroundPath
	case keywordOptionGainLevel:
		gainLevel, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected gain level: %s", ln)
		}

		switch gainLevel {
		case keywordOptionGainLevelVeryLow:
			options.GainLevel = t.GainLevelVeryLow
		case keywordOptionGainLevelLow:
			options.GainLevel = t.GainLevelLow
		case keywordOptionGainLevelMedium:
			options.GainLevel = t.GainLevelMedium
		case keywordOptionGainLevelHigh:
			options.GainLevel = t.GainLevelHigh
		case keywordOptionGainLevelVeryHigh:
			options.GainLevel = t.GainLevelVeryHigh
		default:
			return fmt.Errorf("invalid gain level: %q", gainLevel)
		}
	default:
		return fmt.Errorf("invalid option: %q", option)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return fmt.Errorf("unexpected token after option definition: %q", unknown)
	}

	return nil
}
