package parser

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	keywordComment                 = "#"          // Represents a comment
	keywordOption                  = "@"          // Represents an option
	keywordOptionSampleRate        = "samplerate" // Represents a sample rate option
	keywordOptionVolume            = "volume"     // Represents a volume option
	keywordOptionBackground        = "background" // Represents a background option
	keywordOptionGainLevel         = "gainlevel"  // Represents a gain level option
	keywordOptionGainLevelVeryLow  = "verylow"    // Represents a very low gain level option
	keywordOptionGainLevelLow      = "low"        // Represents a low gain level option
	keywordOptionGainLevelMedium   = "medium"     // Represents a medium gain level option
	keywordOptionGainLevelHigh     = "high"       // Represents a high gain level option
	keywordOptionGainLevelVeryHigh = "veryhigh"   // Represents a very high gain level option
	keywordWaveform                = "waveform"   // Represents a waveform
	keywordSine                    = "sine"       // Represents a sine wave
	keywordSquare                  = "square"     // Represents a square wave
	keywordTriangle                = "triangle"   // Represents a triangle wave
	keywordSawtooth                = "sawtooth"   // Represents a sawtooth wave
	keywordTone                    = "tone"       // Represents a tone
	keywordBinaural                = "binaural"   // Represents a binaural tone
	keywordMonaural                = "monaural"   // Represents a monaural tone
	keywordIsochronic              = "isochronic" // Represents an isochronic tone
	keywordAmplitude               = "amplitude"  // Represents an amplitude
	keywordNoise                   = "noise"      // Represents a noise
	keywordWhite                   = "white"      // Represents a white noise
	keywordPink                    = "pink"       // Represents a pink noise
	keywordBrown                   = "brown"      // Represents a brown noise
	keywordSpin                    = "spin"       // Represents a spin
	keywordWidth                   = "width"      // Represents a width
	keywordRate                    = "rate"       // Represents a rate
	keywordEffect                  = "effect"     // Represents an effect
	keywordBackground              = "background" // Represents a background
	keywordPulse                   = "pulse"      // Represents a pulse
	keywordIntensity               = "intensity"  // Represents an intensity
)

// ParserContext holds the context for parsing
type ParserContext struct {
	Line lineContext // Context for the current line
}

// lineContext holds the context for the current line being parsed
type lineContext struct {
	Raw    string   // Raw line text
	Tokens []string // Tokens extracted from the line
	tkIdx  int      // Current token index
}

// Peek retrieves the next token without advancing the index
func (ctx *lineContext) Peek() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		return ctx.Tokens[ctx.tkIdx], true
	}
	return "", false
}

// NextToken retrieves the next token from the context
func (ctx *lineContext) NextToken() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		token := ctx.Tokens[ctx.tkIdx]
		ctx.tkIdx++
		return token, true
	}
	return "", false
}

// RewindToken moves the token index back by n positions
func (ctx *lineContext) RewindToken(n int) {
	if n <= 0 {
		return
	}
	if n > ctx.tkIdx {
		ctx.tkIdx = 0
		return
	}
	ctx.tkIdx -= n
}

// NextExpectOneOf checks if the next token is one of the expected values
func (ctx *lineContext) NextExpectOneOf(wants ...string) (string, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return "", fmt.Errorf("expected one of %v, got EOF: %s", wants, ctx.Raw)
	}

	if slices.Contains(wants, tok) {
		return tok, nil
	}
	return "", fmt.Errorf("expected one of %v, got %q: %s", wants, tok, ctx.Raw)
}

// NextFloat64Strict retrieves the next token as a float64, enforcing strict parsing
func (ctx *lineContext) NextFloat64Strict() (float64, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, fmt.Errorf("expected float, got EOF: %s", ctx.Raw)
	}

	f, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %q", tok)
	}
	return f, nil
}

// NextIntStrict retrieves the next token as an int, enforcing strict parsing
func (ctx *lineContext) NextIntStrict() (int, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, fmt.Errorf("expected integer, got EOF: %s", ctx.Raw)
	}

	i, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %q", tok)
	}
	return i, nil
}

// NewParserContext creates a new ParserContext for the given line
func NewParserContext(line string) *ParserContext {
	return &ParserContext{
		Line: lineContext{
			Raw:    line,
			Tokens: strings.Fields(line),
			tkIdx:  0,
		},
	}
}
