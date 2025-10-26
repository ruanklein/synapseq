/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"testing"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

func TestHasTrackOverride(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"  track 1 amplitude 10", true},
		{"  track 2 tone 440", true},
		{"  track 3 binaural 8", true},
		{" track 1 amplitude 10", false},               // only 1 space
		{"   track 1 amplitude 10", false},             // 3 spaces
		{"track 1 amplitude 10", false},                // no indentation
		{"  tone 440 binaural 10 amplitude 20", false}, // regular track, not override
		{"  noise pink amplitude 30", false},           // regular track
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasTrackOverride()
		if result != test.expected {
			ts.Errorf("For line %q, expected HasTrackOverride()=%v, got %v",
				test.line, test.expected, result)
		}
	}
}

func TestParseTrackOverride_Success(ts *testing.T) {
	// Create base template preset
	templatePreset, err := t.NewPreset("base", true, nil)
	if err != nil {
		ts.Fatalf("failed to create template: %v", err)
	}

	// Setup template tracks
	templatePreset.Track[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	templatePreset.Track[1] = t.Track{
		Type:      t.TrackBackground,
		Carrier:   200,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(40),
		Waveform:  t.WaveformSine,
		Effect:    t.Effect{Type: t.EffectSpin, Intensity: t.IntensityPercentToRaw(75)},
	}
	templatePreset.Track[2] = t.Track{
		Type:      t.TrackMonauralBeat,
		Carrier:   440,
		Resonance: 8,
		Amplitude: t.AmplitudePercentToRaw(15),
		Waveform:  t.WaveformSquare,
	}

	// Create derived preset
	derivedPreset, err := t.NewPreset("derived", false, templatePreset)
	if err != nil {
		ts.Fatalf("failed to create derived preset: %v", err)
	}

	tests := []struct {
		name      string
		line      string
		trackIdx  int
		checkFunc func(*testing.T, *t.Preset)
	}{
		{
			name:     "override tone carrier",
			line:     "  track 1 tone 350",
			trackIdx: 0,
			checkFunc: func(t *testing.T, p *t.Preset) {
				if p.Track[0].Carrier != 350 {
					t.Errorf("expected carrier 350, got %v", p.Track[0].Carrier)
				}
			},
		},
		{
			name:     "override binaural resonance",
			line:     "  track 1 binaural 12",
			trackIdx: 0,
			checkFunc: func(t *testing.T, p *t.Preset) {
				if p.Track[0].Resonance != 12 {
					t.Errorf("expected resonance 12, got %v", p.Track[0].Resonance)
				}
			},
		},
		{
			name:     "override amplitude",
			line:     "  track 1 amplitude 30",
			trackIdx: 0,
			checkFunc: func(ts *testing.T, p *t.Preset) {
				expected := t.AmplitudePercentToRaw(30)
				if p.Track[0].Amplitude != expected {
					ts.Errorf("expected amplitude %v, got %v", expected, p.Track[0].Amplitude)
				}
			},
		},
		{
			name:     "override background spin width",
			line:     "  track 2 spin 250",
			trackIdx: 1,
			checkFunc: func(t *testing.T, p *t.Preset) {
				if p.Track[1].Carrier != 250 {
					t.Errorf("expected spin width 250, got %v", p.Track[1].Carrier)
				}
			},
		},
		{
			name:     "override background spin rate",
			line:     "  track 2 rate 7",
			trackIdx: 1,
			checkFunc: func(t *testing.T, p *t.Preset) {
				if p.Track[1].Resonance != 7 {
					t.Errorf("expected spin rate 7, got %v", p.Track[1].Resonance)
				}
			},
		},
		{
			name:     "override effect intensity",
			line:     "  track 2 intensity 80",
			trackIdx: 1,
			checkFunc: func(ts *testing.T, p *t.Preset) {
				expected := t.IntensityPercentToRaw(80)
				if p.Track[1].Effect.Intensity != expected {
					ts.Errorf("expected intensity %v, got %v", expected, p.Track[1].Effect.Intensity)
				}
			},
		},
		{
			name:     "override monaural resonance",
			line:     "  track 3 monaural 10",
			trackIdx: 2,
			checkFunc: func(t *testing.T, p *t.Preset) {
				if p.Track[2].Resonance != 10 {
					t.Errorf("expected resonance 10, got %v", p.Track[2].Resonance)
				}
			},
		},
	}

	for _, tt := range tests {
		// Reset derived preset for each test
		derivedPreset.Track = templatePreset.Track

		ctx := NewTextParser(tt.line)
		err := ctx.ParseTrackOverride(derivedPreset)
		if err != nil {
			ts.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}

		tt.checkFunc(ts, derivedPreset)
	}
}

func TestParseTrackOverride_Errors(ts *testing.T) {
	// Create base template preset
	templatePreset, err := t.NewPreset("base", true, nil)
	if err != nil {
		ts.Fatalf("failed to create template: %v", err)
	}

	templatePreset.Track[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	templatePreset.Track[1] = t.Track{
		Type:      t.TrackBackground,
		Carrier:   200,
		Resonance: 5,
		Amplitude: t.AmplitudePercentToRaw(40),
		Waveform:  t.WaveformSine,
		Effect:    t.Effect{Type: t.EffectSpin, Intensity: t.IntensityPercentToRaw(75)},
	}
	templatePreset.Track[2] = t.Track{
		Type:      t.TrackBackground,
		Resonance: 2.5,
		Amplitude: t.AmplitudePercentToRaw(40),
		Waveform:  t.WaveformSine,
		Effect:    t.Effect{Type: t.EffectPulse, Intensity: t.IntensityPercentToRaw(60)},
	}

	// Create derived preset
	derivedPreset, err := t.NewPreset("derived", false, templatePreset)
	if err != nil {
		ts.Fatalf("failed to create derived preset: %v", err)
	}

	tests := []struct {
		name string
		line string
	}{
		{"track index out of range (too high)", "  track 20 amplitude 10"},
		{"track index out of range (zero)", "  track 0 amplitude 10"},
		{"track index out of range (negative)", "  track -1 amplitude 10"},
		{"override off track", "  track 5 amplitude 10"},
		{"missing track index", "  track amplitude 10"},
		{"invalid track index", "  track abc amplitude 10"},
		{"missing value", "  track 1 amplitude"},
		{"invalid value", "  track 1 amplitude abc"},
		{"extra tokens", "  track 1 amplitude 10 extra"},
		{"tone on background track", "  track 2 tone 300"},
		{"spin on non-background track", "  track 1 spin 200"},
		{"wrong beat type override", "  track 1 monaural 8"},
		{"rate on pulse effect", "  track 3 rate 5"},
		{"pulse on spin effect", "  track 2 pulse 3"},
		{"invalid amplitude (too high)", "  track 1 amplitude 150"},
		{"invalid intensity (too high)", "  track 2 intensity 150"},
	}

	for _, tt := range tests {
		// Reset derived preset for each test
		derivedPreset.Track = templatePreset.Track

		ctx := NewTextParser(tt.line)
		err := ctx.ParseTrackOverride(derivedPreset)
		if err == nil {
			ts.Errorf("%s: expected error but got none for line: %q", tt.name, tt.line)
		}
	}
}

func TestParseTrackOverride_WithoutFromPreset(ts *testing.T) {
	// Create preset without 'from' source
	preset, err := t.NewPreset("standalone", false, nil)
	if err != nil {
		ts.Fatalf("failed to create preset: %v", err)
	}

	ctx := NewTextParser("  track 1 amplitude 10")
	err = ctx.ParseTrackOverride(preset)
	if err == nil {
		ts.Errorf("expected error when overriding track on preset without 'from' source")
	}
}

func TestPresetInheritance_Integration(ts *testing.T) {
	// This test simulates the full workflow: template creation, inheritance, and override
	var presets []t.Preset

	// Step 1: Create template preset
	templateCtx := NewTextParser("base-template as template")
	templatePreset, err := templateCtx.ParsePreset(&presets)
	if err != nil {
		ts.Fatalf("failed to create template: %v", err)
	}

	// Add tracks to template
	templatePreset.Track[0] = t.Track{
		Type:      t.TrackBinauralBeat,
		Carrier:   300,
		Resonance: 10,
		Amplitude: t.AmplitudePercentToRaw(20),
		Waveform:  t.WaveformSine,
	}
	templatePreset.Track[1] = t.Track{
		Type:      t.TrackPinkNoise,
		Amplitude: t.AmplitudePercentToRaw(30),
	}

	presets = append(presets, *templatePreset)

	// Step 2: Create derived preset
	derivedCtx := NewTextParser("alpha from base-template")
	derivedPreset, err := derivedCtx.ParsePreset(&presets)
	if err != nil {
		ts.Fatalf("failed to create derived preset: %v", err)
	}

	// Verify inheritance
	if derivedPreset.From == nil {
		ts.Fatalf("derived preset should have 'from' reference")
	}
	if derivedPreset.Track[0].Carrier != 300 {
		ts.Errorf("expected inherited carrier 300, got %v", derivedPreset.Track[0].Carrier)
	}

	// Step 3: Apply overrides
	overrides := []string{
		"  track 1 tone 350",
		"  track 1 binaural 12",
		"  track 1 amplitude 25",
		"  track 2 amplitude 35",
	}

	for _, line := range overrides {
		ctx := NewTextParser(line)
		if err := ctx.ParseTrackOverride(derivedPreset); err != nil {
			ts.Fatalf("failed to apply override %q: %v", line, err)
		}
	}

	// Verify overrides
	if derivedPreset.Track[0].Carrier != 350 {
		ts.Errorf("expected overridden carrier 350, got %v", derivedPreset.Track[0].Carrier)
	}
	if derivedPreset.Track[0].Resonance != 12 {
		ts.Errorf("expected overridden resonance 12, got %v", derivedPreset.Track[0].Resonance)
	}
	if derivedPreset.Track[0].Amplitude != t.AmplitudePercentToRaw(25) {
		ts.Errorf("expected overridden amplitude for track 1, got %v", derivedPreset.Track[0].Amplitude)
	}
	if derivedPreset.Track[1].Amplitude != t.AmplitudePercentToRaw(35) {
		ts.Errorf("expected overridden amplitude for track 2, got %v", derivedPreset.Track[1].Amplitude)
	}

	// Verify that template tracks are unchanged
	if templatePreset.Track[0].Carrier != 300 {
		ts.Errorf("template should remain unchanged, expected carrier 300, got %v", templatePreset.Track[0].Carrier)
	}
}
