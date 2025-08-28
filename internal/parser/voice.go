package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// HasVoice checks if the current line is a voice definition
func (ctx *TextParser) HasVoice() bool {
	ln := ctx.Line.Raw

	if len(ln) < 3 {
		return false
	}

	return ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' '
}

// ParseVoice extracts and returns a Voice from the current line context
func (ctx *TextParser) ParseVoice() (*t.Voice, error) {
	waveform := t.WaveformSine
	ln := ctx.Line.Raw

	if tok, ok := ctx.Line.Peek(); ok && tok == t.KeywordWaveform {
		ctx.Line.NextToken() // skip "waveform"

		wfTok, err := ctx.Line.NextExpectOneOf(t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, %q, or %q after waveform: %s", t.KeywordSine, t.KeywordSquare, t.KeywordTriangle, t.KeywordSawtooth, ln)
		}

		switch wfTok {
		case t.KeywordSine:
			waveform = t.WaveformSine
		case t.KeywordSquare:
			waveform = t.WaveformSquare
		case t.KeywordTriangle:
			waveform = t.WaveformTriangle
		case t.KeywordSawtooth:
			waveform = t.WaveformSawtooth
		}

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordSpin, t.KeywordEffect); err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after waveform type: %s", t.KeywordTone, t.KeywordSpin, t.KeywordEffect, ln)
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q: %s", t.KeywordTone, t.KeywordNoise, t.KeywordSpin, t.KeywordEffect, t.KeywordBackground, ln)
	}

	var (
		carrier, resonance, amplitude, intensity float64
		voiceType                                t.VoiceType
	)

	switch first {
	case t.KeywordTone:
		var err error
		if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}

		kind, err := ctx.Line.NextExpectOneOf(t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after carrier: %s", t.KeywordBinaural, t.KeywordMonaural, t.KeywordIsochronic, ln)
		}

		switch kind {
		case t.KeywordBinaural:
			voiceType = t.VoiceBinauralBeat
		case t.KeywordMonaural:
			voiceType = t.VoiceMonauralBeat
		case t.KeywordIsochronic:
			voiceType = t.VoiceIsochronicBeat
		}

		if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case t.KeywordNoise:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordWhite, t.KeywordPink, t.KeywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after noise: %s", t.KeywordWhite, t.KeywordPink, t.KeywordBrown, ln)
		}

		switch kind {
		case t.KeywordWhite:
			voiceType = t.VoiceWhiteNoise
		case t.KeywordPink:
			voiceType = t.VoicePinkNoise
		case t.KeywordBrown:
			voiceType = t.VoiceBrownNoise
		}

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after noise type: %s", t.KeywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case t.KeywordSpin:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordWhite, t.KeywordPink, t.KeywordBrown)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q, or %q after spin: %s", t.KeywordWhite, t.KeywordPink, t.KeywordBrown, ln)
		}

		switch kind {
		case t.KeywordWhite:
			voiceType = t.VoiceSpinWhite
		case t.KeywordPink:
			voiceType = t.VoiceSpinPink
		case t.KeywordBrown:
			voiceType = t.VoiceSpinBrown
		}

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordWidth); err != nil {
			return nil, fmt.Errorf("expected %q after spin noise type: %s", t.KeywordWidth, ln)
		}
		if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("carrier: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(t.KeywordRate); err != nil {
			return nil, fmt.Errorf("expected %q after carrier: %s", t.KeywordRate, ln)
		}
		if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("resonance: %w", err)
		}
		if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case t.KeywordBackground:
		voiceType = t.VoiceBackground
		var err error
		if _, err = ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after background: %s", t.KeywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case t.KeywordEffect:
		var err error
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordSpin, t.KeywordPulse)
		if err != nil {
			return nil, fmt.Errorf("expected %q or %q after effect: %s", t.KeywordSpin, t.KeywordPulse, ln)
		}

		switch kind {
		case t.KeywordSpin:
			voiceType = t.VoiceEffectSpin
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordWidth); err != nil {
				return nil, fmt.Errorf("expected %q after spin: %s", t.KeywordWidth, ln)
			}
			if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("carrier: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordRate); err != nil {
				return nil, fmt.Errorf("expected %q after carrier: %s", t.KeywordRate, ln)
			}
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		case t.KeywordPulse:
			voiceType = t.VoiceEffectPulse
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err = ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
		}
	default:
		return nil, fmt.Errorf("expected %q, %q, %q, %q, or %q. Received: %s", t.KeywordTone, t.KeywordNoise, t.KeywordSpin, t.KeywordEffect, t.KeywordBackground, first)
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
