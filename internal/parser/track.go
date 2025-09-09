package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// HasTrack checks if the current line is a track definition
func (ctx *TextParser) HasTrack() bool {
	ln := ctx.Line.Raw

	if len(ln) < 3 {
		return false
	}

	return ln[0] == ' ' && ln[1] == ' ' && ln[2] != ' '
}

// ParseTrack extracts and returns a Track from the current line context
func (ctx *TextParser) ParseTrack() (*t.Track, error) {
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

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordTone, t.KeywordBackground); err != nil {
			return nil, fmt.Errorf("expected %q or %q after waveform type: %s", t.KeywordTone, t.KeywordBackground, ln)
		}

		ctx.Line.RewindToken(1) // rewind to re-process the tone line
	}

	first, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected %q, %s or %q: %s", t.KeywordTone, t.KeywordNoise, t.KeywordBackground, ln)
	}

	var (
		carrier, resonance, amplitude float64
		trackType                     t.TrackType
	)

	effect := t.Effect{Type: t.EffectOff, Intensity: 0.0}

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
			trackType = t.TrackBinauralBeat
		case t.KeywordMonaural:
			trackType = t.TrackMonauralBeat
		case t.KeywordIsochronic:
			trackType = t.TrackIsochronicBeat
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
			trackType = t.TrackWhiteNoise
		case t.KeywordPink:
			trackType = t.TrackPinkNoise
		case t.KeywordBrown:
			trackType = t.TrackBrownNoise
		}

		if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
			return nil, fmt.Errorf("expected %q after noise type: %s", t.KeywordAmplitude, ln)
		}
		if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
			return nil, fmt.Errorf("amplitude: %w", err)
		}
	case t.KeywordBackground:
		trackType = t.TrackBackground
		kind, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude, t.KeywordSpin, t.KeywordPulse)
		if err != nil {
			return nil, fmt.Errorf("expected %q, %q or %q after background: %s", t.KeywordAmplitude, t.KeywordSpin, t.KeywordPulse, ln)
		}

		intensity := 0.0

		switch kind {
		case t.KeywordAmplitude:
			if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("amplitude: %w", err)
			}
		case t.KeywordSpin:
			effect.Type = t.EffectSpin
			if carrier, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("carrier: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordRate); err != nil {
				return nil, fmt.Errorf("expected %q after carrier: %s", t.KeywordRate, ln)
			}
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordAmplitude, ln)
			}
			if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("amplitude: %w", err)
			}
		case t.KeywordPulse:
			effect.Type = t.EffectPulse
			if resonance, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("resonance: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordIntensity); err != nil {
				return nil, fmt.Errorf("expected %q after resonance: %s", t.KeywordIntensity, ln)
			}
			if intensity, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("intensity: %w", err)
			}
			if _, err := ctx.Line.NextExpectOneOf(t.KeywordAmplitude); err != nil {
				return nil, fmt.Errorf("expected %q after intensity: %s", t.KeywordAmplitude, ln)
			}
			if amplitude, err = ctx.Line.NextFloat64Strict(); err != nil {
				return nil, fmt.Errorf("amplitude: %w", err)
			}
		}

		effect.Intensity = t.IntensityPercentToRaw(intensity)
	default:
		return nil, fmt.Errorf("expected %q, %q or %q. Received: %s", t.KeywordTone, t.KeywordNoise, t.KeywordBackground, first)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after track definition: %q", unknown)
	}

	track := t.Track{
		Type:      trackType,
		Carrier:   carrier,
		Resonance: resonance,
		Amplitude: t.AmplitudePercentToRaw(amplitude),
		Waveform:  waveform,
		Effect:    effect,
	}
	if err := track.Validate(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &track, nil
}
