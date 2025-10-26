/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"fmt"
	"strings"
	"testing"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

func TestHasTrack(ts *testing.T) {
	trLnTone := (&t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   440,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(4),
	}).String()

	trLnNoise := (&t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}).String()

	trLnBackground := (&t.Track{
		Type:      t.TrackBackground,
		Amplitude: t.AmplitudePercentToRaw(50),
	}).String()

	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("  %s", trLnTone), true},
		{fmt.Sprintf("  %s", trLnNoise), true},
		{fmt.Sprintf("  %s", trLnBackground), true},
		{fmt.Sprintf(" %s", trLnTone), false},
		{fmt.Sprintf("   %s", trLnTone), false},
		{trLnTone, false},
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
	trs := []*t.Track{
		{
			Type:      t.TrackBinauralBeat,
			Carrier:   300,
			Resonance: 10,
			Amplitude: t.AmplitudePercentToRaw(15),
		},
		{
			Type:      t.TrackMonauralBeat,
			Carrier:   440,
			Resonance: 11,
			Amplitude: t.AmplitudePercentToRaw(20),
		},
		{
			Type:      t.TrackIsochronicBeat,
			Carrier:   220,
			Resonance: 8,
			Amplitude: t.AmplitudePercentToRaw(5),
		},
		{
			Type:      t.TrackBinauralBeat,
			Carrier:   300,
			Resonance: 10,
			Amplitude: t.AmplitudePercentToRaw(15),
			Waveform:  t.WaveformTriangle,
		},
		{
			Type:      t.TrackPureTone,
			Carrier:   350,
			Amplitude: t.AmplitudePercentToRaw(10),
			Waveform:  t.WaveformSquare,
		},
	}

	// Helper to format track without extra waveform
	fmtLine := func(tr *t.Track) string {
		return strings.Join(strings.Fields(tr.String())[2:], " ")
	}

	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{fmtLine(trs[0]), *trs[0]},
		{fmtLine(trs[1]), *trs[1]},
		{fmtLine(trs[2]), *trs[2]},
		{trs[3].String(), *trs[3]},
		{trs[4].String(), *trs[4]},
	}

	for i, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("Test %d: For line '%s', expected track %+v but got %+v", i, tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Noise(ts *testing.T) {
	trs := []*t.Track{
		{
			Type:      t.TrackWhiteNoise,
			Amplitude: t.AmplitudePercentToRaw(5),
		},
		{
			Type:      t.TrackPinkNoise,
			Amplitude: t.AmplitudePercentToRaw(40),
		},
		{
			Type:      t.TrackBrownNoise,
			Amplitude: t.AmplitudePercentToRaw(15),
		},
	}

	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{trs[0].String(), *trs[0]},
		{trs[1].String(), *trs[1]},
		{trs[2].String(), *trs[2]},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
		}
	}
}

func TestParseTrack_Background(ts *testing.T) {
	trs := []*t.Track{
		{
			Type:      t.TrackBackground,
			Amplitude: t.AmplitudePercentToRaw(50),
		},
		{
			Type:      t.TrackBackground,
			Carrier:   200,
			Resonance: 5,
			Effect:    t.Effect{Type: t.EffectSpin, Intensity: t.IntensityPercentToRaw(75)},
			Amplitude: t.AmplitudePercentToRaw(50),
		},
		{
			Type:      t.TrackBackground,
			Resonance: 2.5,
			Effect:    t.Effect{Type: t.EffectPulse, Intensity: t.IntensityPercentToRaw(60)},
			Amplitude: t.AmplitudePercentToRaw(40),
		},
		{
			Type:      t.TrackBackground,
			Resonance: 2.5,
			Effect:    t.Effect{Type: t.EffectPulse, Intensity: t.IntensityPercentToRaw(60)},
			Amplitude: t.AmplitudePercentToRaw(40),
			Waveform:  t.WaveformSquare,
		},
		{
			Type:      t.TrackBackground,
			Amplitude: t.AmplitudePercentToRaw(33),
		},
	}

	tests := []struct {
		line      string
		wantTrack t.Track
	}{
		{trs[0].String(), *trs[0]},
		{trs[1].String(), *trs[1]},
		{trs[2].String(), *trs[2]},
		{trs[3].String(), *trs[3]},
		{trs[4].String(), *trs[4]},
	}

	for _, tt := range tests {
		ctx := NewTextParser(tt.line)
		tr, err := ctx.ParseTrack()
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", tt.line, err)
			continue
		}
		if *tr != tt.wantTrack {
			ts.Errorf("For line '%s', expected track %+v but got %+v", tt.line, tt.wantTrack, *tr)
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
