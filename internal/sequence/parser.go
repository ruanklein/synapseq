package sequence

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ruanklein/synapseq/internal/audio"
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

// isOptionLine checks if the first element is an option
func isOptionLine(ctx *LineContext) bool {
	if len(ctx.Line) == 0 {
		return false
	}
	return string(ctx.Line[0]) == keywordOption
}

// parseOption extracts and applies the option from the elements
func parseOptionLine(ctx *LineContext, options *SequenceOptions) error {
	tok, ok := ctx.NextToken()
	if !ok {
		return fmt.Errorf("expected option, got EOF: %s", ctx.Line)
	}
	if string(tok[0]) != keywordOption {
		return fmt.Errorf("expected option. Received: %s", tok)
	}

	option := tok[1:]
	if len(option) == 0 {
		return fmt.Errorf("expected option name: %s", ctx.Line)
	}

	switch option {
	case keywordOptionSampleRate:
		sampleRate, err := ctx.NextIntStrict()
		if err != nil {
			return fmt.Errorf("samplerate: %v", err)
		}
		options.SampleRate = sampleRate
	case keywordOptionVolume:
		volume, err := ctx.NextIntStrict()
		if err != nil {
			return fmt.Errorf("volume: %v", err)
		}
		options.Volume = volume
	case keywordOptionBackground:
		backgroundPath, ok := ctx.NextToken()
		if !ok {
			return fmt.Errorf("expected background path: %s", ctx.Line)
		}
		options.BackgroundPath = backgroundPath
	case keywordOptionGainLevel:
		gainLevel, ok := ctx.NextToken()
		if !ok {
			return fmt.Errorf("expected gain level: %s", ctx.Line)
		}

		switch gainLevel {
		case keywordOptionGainLevelVeryLow:
			options.GainLevel = gainVeryLow
		case keywordOptionGainLevelLow:
			options.GainLevel = gainLow
		case keywordOptionGainLevelMedium:
			options.GainLevel = gainMedium
		case keywordOptionGainLevelHigh:
			options.GainLevel = gainHigh
		case keywordOptionGainLevelVeryHigh:
			options.GainLevel = gainVeryHigh
		default:
			return fmt.Errorf("invalid gain level: %q", gainLevel)
		}
	default:
		return fmt.Errorf("invalid option: %q", option)
	}

	// Check for unexpected tokens
	unknown, ok := ctx.NextToken()
	if ok {
		return fmt.Errorf("unexpected token after definition: %q", unknown)
	}

	return nil
}

// isCommentLine checks if the first element is a comment
func isCommentLine(ctx *LineContext) bool {
	tok, ok := ctx.Peek()
	return ok && string(tok[0]) == keywordComment
}

// parseComment extracts and prints the comment from the elements
func parseCommentLine(ctx *LineContext) string {
	tok, ok := ctx.Peek()
	if !ok || string(tok[0]) != keywordComment {
		return ""
	}
	if len(tok) >= 2 && string(tok[1]) == keywordComment {
		return fmt.Sprintf("%s %s", tok[2:], strings.Join(ctx.Tokens[1:], " "))
	}
	return ""
}

// isValidPresetName checks if a string is a valid preset name
func isValidPresetName(s string) bool {
	regexPreset := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	return regexPreset.MatchString(s)
}

// parsePreset extracts and returns a Preset from the current line context
func parsePresetLine(ctx *LineContext) (*Preset, error) {
	tok, ok := ctx.Peek()
	if !ok {
		return nil, fmt.Errorf("expected preset name, got EOF: %s", ctx.Line)
	}

	presetName := strings.ToLower(tok)
	if presetName == builtinSilence {
		return nil, fmt.Errorf("cannot load %q built-in preset: %s", presetName, ctx.Line)
	}

	preset := &Preset{Name: presetName}
	preset.InitVoices()
	return preset, nil
}

// isVoiceLine checks if the current line is a voice definition
func isVoiceLine(ctx *LineContext) bool {
	if len(ctx.Line) < 3 {
		return false
	}
	return ctx.Line[0] == ' ' && ctx.Line[1] == ' ' && ctx.Line[2] != ' '
}

// parseVoiceLine extracts and returns a Voice from the current line context
func parseVoiceLine(ctx *LineContext) (*audio.Voice, error) {
	waveform := audio.WaveformSine

	if tok, ok := ctx.Peek(); ok && tok == keywordWaveform {
		ctx.NextToken() // skip "waveform"

		wfTok, err := ctx.NextExpectOneOf(keywordSine, keywordSquare, keywordTriangle, keywordSawtooth)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, %q, or %q after waveform: %s", keywordSine, keywordSquare, keywordTriangle, keywordSawtooth, ctx.Line)
		}

		switch wfTok {
		case keywordSine:
			waveform = audio.WaveformSine
		case keywordSquare:
			waveform = audio.WaveformSquare
		case keywordTriangle:
			waveform = audio.WaveformTriangle
		case keywordSawtooth:
			waveform = audio.WaveformSawtooth
		}

		if _, err := ctx.NextExpectOneOf(keywordTone, keywordSpin, keywordEffect); err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after waveform type: %s", keywordTone, keywordSpin, keywordEffect, ctx.Line)
		}

		ctx.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q: %s", keywordTone, keywordNoise, keywordSpin, keywordEffect, keywordBackground, ctx.Line)
	}

	var (
		carrier, resonance, amplitude, intensity float64
		voiceType                                audio.VoiceType
	)

	switch first {
	case keywordTone:
		var err error
		if carrier, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}

		kind, err := ctx.NextExpectOneOf(keywordBinaural, keywordMonaural, keywordIsochronic)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after carrier: %s", keywordBinaural, keywordMonaural, keywordIsochronic, ctx.Line)
		}

		switch kind {
		case keywordBinaural:
			voiceType = audio.VoiceBinauralBeat
		case keywordMonaural:
			voiceType = audio.VoiceMonauralBeat
		case keywordIsochronic:
			voiceType = audio.VoiceIsochronicBeat
		}

		if resonance, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", keywordAmplitude, ctx.Line)
		}
		if amplitude, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordNoise:
		var err error
		kind, err := ctx.NextExpectOneOf(keywordWhite, keywordPink, keywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after noise: %s", keywordWhite, keywordPink, keywordBrown, ctx.Line)
		}

		switch kind {
		case keywordWhite:
			voiceType = audio.VoiceWhiteNoise
		case keywordPink:
			voiceType = audio.VoicePinkNoise
		case keywordBrown:
			voiceType = audio.VoiceBrownNoise
		}

		if _, err := ctx.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after noise type: %s", keywordAmplitude, ctx.Line)
		}
		if amplitude, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordSpin:
		var err error
		kind, err := ctx.NextExpectOneOf(keywordWhite, keywordPink, keywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after spin: %s", keywordWhite, keywordPink, keywordBrown, ctx.Line)
		}

		switch kind {
		case keywordWhite:
			voiceType = audio.VoiceSpinWhite
		case keywordPink:
			voiceType = audio.VoiceSpinPink
		case keywordBrown:
			voiceType = audio.VoiceSpinBrown
		}

		if _, err := ctx.NextExpectOneOf(keywordWidth); err != nil {
			return nil, fmt.Errorf("expected %q after spin noise type: %s", keywordWidth, ctx.Line)
		}
		if carrier, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}
		if _, err := ctx.NextExpectOneOf(keywordRate); err != nil {
			return nil, fmt.Errorf("expected %q after carrier: %s", keywordRate, ctx.Line)
		}
		if resonance, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", keywordAmplitude, ctx.Line)
		}
		if amplitude, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordBackground:
		voiceType = audio.VoiceBackground
		var err error
		if _, err = ctx.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after background: %s", keywordAmplitude, ctx.Line)
		}
		if amplitude, err = ctx.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordEffect:
		var err error
		kind, err := ctx.NextExpectOneOf(keywordSpin, keywordPulse)
		if err != nil {
			return nil, fmt.Errorf("expected %q or %q after effect: %s", keywordSpin, keywordPulse, ctx.Line)
		}

		switch kind {
		case keywordSpin:
			voiceType = audio.VoiceEffectSpin
			if _, err := ctx.NextExpectOneOf(keywordWidth); err != nil {
				return nil, fmt.Errorf("expected %q after spin: %s", keywordWidth, ctx.Line)
			}
			if carrier, err = ctx.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("carrier: %w", err)
			}
			if _, err := ctx.NextExpectOneOf(keywordRate); err != nil {
				return nil, fmt.Errorf("expected %q after carrier: %s", keywordRate, ctx.Line)
			}
			if resonance, err = ctx.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.NextExpectOneOf(keywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", keywordIntensity, ctx.Line)
			}
			if intensity, err = ctx.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		case keywordPulse:
			voiceType = audio.VoiceEffectPulse
			if resonance, err = ctx.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.NextExpectOneOf(keywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", keywordIntensity, ctx.Line)
			}
			if intensity, err = ctx.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q. Received: %s", keywordTone, keywordNoise, keywordSpin, keywordEffect, keywordBackground, first)
	}

	// Check for unexpected tokens
	unknown, ok := ctx.NextToken()
	if ok {
		return nil, fmt.Errorf("unexpected token after definition: %q", unknown)
	}

	voice := audio.Voice{
		Type:      voiceType,
		Carrier:   carrier,
		Resonance: resonance,
		Amplitude: audio.AmplitudePercentToRaw(amplitude),
		Intensity: audio.IntensityPercentToRaw(intensity),
		Waveform:  waveform,
	}
	if err := voice.Validate(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &voice, nil
}
