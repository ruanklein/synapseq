package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// IsOption checks if the first element is an option
func (ctx *TextParser) IsOption() bool {
	ln := ctx.Line.Raw

	if len(ln) == 0 {
		return false
	}

	return string(ln[0]) == t.KeywordOption
}

// ParseOption extracts and applies the option from the elements
func (ctx *TextParser) ParseOption(options *t.Option) error {
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
	case t.KeywordOptionBackground:
		backgroundPath, ok := ctx.Line.NextToken()
		if !ok {
			return fmt.Errorf("expected background path: %s", ln)
		}
		options.BackgroundPath = backgroundPath
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

	// Validate options
	if err := options.Validate(); err != nil {
		return fmt.Errorf("%v", err)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return fmt.Errorf("unexpected token after option definition: %q", unknown)
	}

	return nil
}
