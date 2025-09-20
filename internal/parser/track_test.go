/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"testing"

	t "github.com/ruanklein/synapseq/internal/types"
)

func TestHasTrack(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"  tone 300 binaural 10 amplitude 10", true},
		{"  noise pink amplitude 30", true},
		{"  background amplitude 50", true},
		{" tone 300 binaural 10 amplitude 10", false},
		{"   tone 300 binaural 10 amplitude 10", false},
		{"tone 300 binaural 10 amplitude 10", false},
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		got := ctx.HasTrack()
		if got != test.expected {
			ts.Errorf("For line '%s', expected HasTrack()=%v, got %v", test.line, test.expected, got)
		}
	}
}

func TestParseTrack_Tones(ts *testing.T) {
	tests := []struct {
		line       string
		wantType   t.TrackType
		wantWF     t.WaveformType
		carrier    float64
		resonance  float64
		amplitudeP float64
	}{
		{"  tone 300 binaural 10 amplitude 15", t.TrackBinauralBeat, t.WaveformSine, 300, 10, 15},
		{"  tone 440 monaural 11 amplitude 20", t.TrackMonauralBeat, t.WaveformSine, 440, 11, 20},
		{"  tone 220 isochronic 8 amplitude 5", t.TrackIsochronicBeat, t.WaveformSine, 220, 8, 5},
		{"  waveform triangle tone 300 binaural 10 amplitude 15", t.TrackBinauralBeat, t.WaveformTriangle, 300, 10, 15},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if tr.Type != tt.wantType {
			ts.Errorf("For line '%s', want Type %v, got %v", tt.line, tt.wantType, tr.Type)
		}
		if tr.Waveform != tt.wantWF {
			ts.Errorf("For line '%s', want Waveform %v, got %v", tt.line, tt.wantWF, tr.Waveform)
		}
		if tr.Carrier != tt.carrier {
			ts.Errorf("For line '%s', want Carrier %.4f, got %.4f", tt.line, tt.carrier, tr.Carrier)
		}
		if tr.Resonance != tt.resonance {
			ts.Errorf("For line '%s', want Resonance %.4f, got %.4f", tt.line, tt.resonance, tr.Resonance)
		}
		if tr.Effect.Type != t.EffectOff {
			ts.Errorf("For line '%s', want Effect Off, got %v", tt.line, tr.Effect.Type)
		}
		if tr.Effect.Intensity != t.IntensityPercentToRaw(0) {
			ts.Errorf("For line '%s', want Intensity 0, got %.4f", tt.line, tr.Effect.Intensity)
		}
		if tr.Amplitude != t.AmplitudePercentToRaw(tt.amplitudeP) {
			ts.Errorf("For line '%s', want Amplitude raw %.4f, got %.4f", tt.line, t.AmplitudePercentToRaw(tt.amplitudeP), tr.Amplitude)
		}
	}
}

func TestParseTrack_Noise(ts *testing.T) {
	tests := []struct {
		line       string
		wantType   t.TrackType
		amplitudeP float64
	}{
		{"  noise white amplitude 5", t.TrackWhiteNoise, 5},
		{"  noise pink amplitude 40", t.TrackPinkNoise, 40},
		{"  noise brown amplitude 15", t.TrackBrownNoise, 15},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if tr.Type != tt.wantType {
			ts.Errorf("For line '%s', want Type %v, got %v", tt.line, tt.wantType, tr.Type)
		}
		if tr.Waveform != t.WaveformSine {
			ts.Errorf("For line '%s', want default waveform sine, got %v", tt.line, tr.Waveform)
		}
		if tr.Amplitude != t.AmplitudePercentToRaw(tt.amplitudeP) {
			ts.Errorf("For line '%s', want Amplitude raw %.4f, got %.4f", tt.line, t.AmplitudePercentToRaw(tt.amplitudeP), tr.Amplitude)
		}
	}
}

func TestParseTrack_Background(ts *testing.T) {
	tests := []struct {
		line       string
		wantWF     t.WaveformType
		wantType   t.TrackType
		wantEff    t.EffectType
		carrier    float64
		resonance  float64
		intensityP float64
		amplitudeP float64
	}{
		{"  background amplitude 50", t.WaveformSine, t.TrackBackground, t.EffectOff, 0, 0, 0, 50},
		{"  background spin 200 rate 5 intensity 75 amplitude 50", t.WaveformSine, t.TrackBackground, t.EffectSpin, 200, 5, 75, 50},
		{"  background pulse 2.5 intensity 60 amplitude 40", t.WaveformSine, t.TrackBackground, t.EffectPulse, 0, 2.5, 60, 40},
		{"  waveform sawtooth background amplitude 33", t.WaveformSawtooth, t.TrackBackground, t.EffectOff, 0, 0, 0, 33},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if tr.Type != tt.wantType {
			ts.Errorf("For line '%s', want Type %v, got %v", tt.line, tt.wantType, tr.Type)
		}
		if tr.Waveform != tt.wantWF {
			ts.Errorf("For line '%s', want Waveform %v, got %v", tt.line, tt.wantWF, tr.Waveform)
		}
		if tr.Effect.Type != tt.wantEff {
			ts.Errorf("For line '%s', want Effect %v, got %v", tt.line, tt.wantEff, tr.Effect.Type)
		}
		if tr.Carrier != tt.carrier {
			ts.Errorf("For line '%s', want Carrier %.4f, got %.4f", tt.line, tt.carrier, tr.Carrier)
		}
		if tr.Resonance != tt.resonance {
			ts.Errorf("For line '%s', want Resonance %.4f, got %.4f", tt.line, tt.resonance, tr.Resonance)
		}
		if tr.Effect.Intensity != t.IntensityPercentToRaw(tt.intensityP) {
			ts.Errorf("For line '%s', want Intensity raw %.4f, got %.4f", tt.line, t.IntensityPercentToRaw(tt.intensityP), tr.Effect.Intensity)
		}
		if tr.Amplitude != t.AmplitudePercentToRaw(tt.amplitudeP) {
			ts.Errorf("For line '%s', want Amplitude raw %.4f, got %.4f", tt.line, t.AmplitudePercentToRaw(tt.amplitudeP), tr.Amplitude)
		}
	}
}

func TestParseTrack_Errors(ts *testing.T) {
	tests := []string{
		"  tone 300 binaural amplitude 10",
		"  tone 300 unknown 10 amplitude 10",
		"  noise white amplitude",
		"  background spin 200 rate five intensity 75 amplitude 50",
		"  background pulse 2.5 intensity sixty amplitude 40",
		"  background amplitude 50 extra",
		"  tone 300 binaural 10 amplitude 120",
		"  background pulse 2.5 intensity 150 amplitude 40",
		"  unknown something",
	}

	for _, line := range tests {
		ctx := NewTextParser(line)
		_, err := ctx.ParseTrack()
		if err == nil {
			ts.Errorf("For line '%s', expected error but got none", line)
		}
	}
}
