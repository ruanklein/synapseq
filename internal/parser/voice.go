package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// IsVoiceLine checks if the current line is a voice definition
func (ctx *ParserContext) IsVoiceLine() bool {
	ln := ctx.Line.Raw

	if len(ln) < 3 {
		return false
	}

	return ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' '
}

// ParseVoiceLine extracts and returns a Voice from the current line context
func (ctx *ParserContext) ParseVoiceLine() (*t.Voice, error) {
	waveform := t.WaveformSine
	ln := ctx.Line.Raw

	if tok, ok := ctx.Line.Peek(); ok && tok == keywordWaveform {
		ctx.Line.NextToken() // skip "waveform"

		wfTok, err := ctx.Line.NextExpectOneOf(keywordSine, keywordSquare, keywordTriangle, keywordSawtooth)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, %q, or %q after waveform: %s", keywordSine, keywordSquare, keywordTriangle, keywordSawtooth, ln)
		}

		switch wfTok {
		case keywordSine:
			waveform = t.WaveformSine
		case keywordSquare:
			waveform = t.WaveformSquare
		case keywordTriangle:
			waveform = t.WaveformTriangle
		case keywordSawtooth:
			waveform = t.WaveformSawtooth
		}

		if _, err := ctx.Line.NextExpectOneOf(keywordTone, keywordSpin, keywordEffect); err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after waveform type: %s", keywordTone, keywordSpin, keywordEffect, ln)
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q: %s", keywordTone, keywordNoise, keywordSpin, keywordEffect, keywordBackground, ln)
	}

	var (
		carrier, resonance, amplitude, intensity float64
		voiceType                                t.VoiceType
	)

	switch first {
	case keywordTone:
		var err error
		if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}

		kind, err := ctx.Line.NextExpectOneOf(keywordBinaural, keywordMonaural, keywordIsochronic)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after carrier: %s", keywordBinaural, keywordMonaural, keywordIsochronic, ln)
		}

		switch kind {
		case keywordBinaural:
			voiceType = t.VoiceBinauralBeat
		case keywordMonaural:
			voiceType = t.VoiceMonauralBeat
		case keywordIsochronic:
			voiceType = t.VoiceIsochronicBeat
		}

		if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", keywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordNoise:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(keywordWhite, keywordPink, keywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after noise: %s", keywordWhite, keywordPink, keywordBrown, ln)
		}

		switch kind {
		case keywordWhite:
			voiceType = t.VoiceWhiteNoise
		case keywordPink:
			voiceType = t.VoicePinkNoise
		case keywordBrown:
			voiceType = t.VoiceBrownNoise
		}

		if _, err := ctx.Line.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after noise type: %s", keywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordSpin:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(keywordWhite, keywordPink, keywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after spin: %s", keywordWhite, keywordPink, keywordBrown, ln)
		}

		switch kind {
		case keywordWhite:
			voiceType = t.VoiceSpinWhite
		case keywordPink:
			voiceType = t.VoiceSpinPink
		case keywordBrown:
			voiceType = t.VoiceSpinBrown
		}

		if _, err := ctx.Line.NextExpectOneOf(keywordWidth); err != nil {
			return nil, fmt.Errorf("expected %q after spin noise type: %s", keywordWidth, ln)
		}
		if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(keywordRate); err != nil {
			return nil, fmt.Errorf("expected %q after carrier: %s", keywordRate, ln)
		}
		if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", keywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordBackground:
		voiceType = t.VoiceBackground
		var err error
		if _, err = ctx.Line.NextExpectOneOf(keywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after background: %s", keywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case keywordEffect:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(keywordSpin, keywordPulse)
		if err != nil {
			return nil, fmt.Errorf("expected %q or %q after effect: %s", keywordSpin, keywordPulse, ln)
		}

		switch kind {
		case keywordSpin:
			voiceType = t.VoiceEffectSpin
			if _, err := ctx.Line.NextExpectOneOf(keywordWidth); err != nil {
				return nil, fmt.Errorf("expected %q after spin: %s", keywordWidth, ln)
			}
			if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("carrier: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(keywordRate); err != nil {
				return nil, fmt.Errorf("expected %q after carrier: %s", keywordRate, ln)
			}
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.Line.NextExpectOneOf(keywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", keywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		case keywordPulse:
			voiceType = t.VoiceEffectPulse
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.Line.NextExpectOneOf(keywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", keywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q. Received: %s", keywordTone, keywordNoise, keywordSpin, keywordEffect, keywordBackground, first)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after voice definition: %q", unknown)
	}

	voice := t.Voice{
		Type:      voiceType,
		Carrier:   carrier,
		Resonance: resonance,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
		Intensity: t.IntensityPercentToRaw(intensity),
		Waveform:  waveform,
	}
	if err := voice.Validate(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &voice, nil
}
