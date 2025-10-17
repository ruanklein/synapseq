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

func TestHasTimeline(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"00:00:00 alpha", true},
		{"23:59:59 preset-1", true},
		{" 00:00:00 alpha", false},
		{"00:00 alpha", false},
		{"alpha", false},
		{"+00:00:10", false},
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasTimeline()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasTimeline() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseTimeline(ts *testing.T) {
	var presets []t.Preset
	alpha, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset 'alpha': %v", err)
	}
	presets = append(presets, *alpha)

	tests := []struct {
		line        string
		expectError bool
		expectedMs  int
	}{
		{"00:00:00 alpha", false, 0},
		{"00:00:15 alpha", false, 15_000},
		{"12:34:56 alpha", false, (12*3600 + 34*60 + 56) * 1000},
		{"24:00:00 alpha", true, 0},
		{"00:60:00 alpha", true, 0},
		{"00:00:60 alpha", true, 0},
		{"00:00:05 beta", true, 0},
		{"00:00:05 alpha extra", true, 0},
		{"00:00:05", true, 0},
		{"00:00 alpha", true, 0},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		per, err := ctx.ParseTimeline(&presets)
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got none", test.line)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}
		if per == nil {
			ts.Errorf("For line '%s', expected non-nil period", test.line)
			continue
		}
		if per.Time != test.expectedMs {
			ts.Errorf("For line '%s', expected time %d but got %d", test.line, test.expectedMs, per.Time)
		}
	}
}

func TestParseTime(ts *testing.T) {
	tests := []struct {
		in          string
		expectedMs  int
		expectError bool
	}{
		{"00:00:00", 0, false},
		{"00:00:01", 1_000, false},
		{"00:01:00", 60_000, false},
		{"01:00:00", 3_600_000, false},
		{"12:34:56", (12*3600 + 34*60 + 56) * 1000, false},
		{"23:59:59", (23*3600 + 59*60 + 59) * 1000, false},

		// Invalid cases
		{"0:00:00", 0, true},
		{"00:0:00", 0, true},
		{"00:00:0", 0, true},
		{"24:00:00", 0, true},
		{"00:60:00", 0, true},
		{"00:00:60", 0, true},
		{"aa:bb:cc", 0, true},
		{"00:00", 0, true},
		{"000000", 0, true},
		{"+00:01:00", 0, true},
		{"", 0, true},
		{"   ", 0, true},
	}

	for _, test := range tests {
		ms, err := parseTime(test.in)
		if test.expectError {
			if err == nil {
				ts.Errorf("For time '%s', expected error but got %d", test.in, ms)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For time '%s', unexpected error: %v", test.in, err)
			continue
		}
		if ms != test.expectedMs {
			ts.Errorf("For time '%s', expected %d but got %d", test.in, test.expectedMs, ms)
		}
	}
}

func TestParseTimelineWithTransitions(ts *testing.T) {
	var presets []t.Preset
	alpha, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating preset 'alpha': %v", err)
	}
	presets = append(presets, *alpha)

	tests := []struct {
		line               string
		expectError        bool
		expectedMs         int
		expectedTransition t.TransitionType
	}{
		// Transition steady (explicit)
		{"00:00:00 alpha steady", false, 0, t.TransitionSteady},
		{"00:00:15 alpha steady", false, 15_000, t.TransitionSteady},

		// Transition ease-out
		{"00:01:00 alpha ease-out", false, 60_000, t.TransitionEaseOut},
		{"12:34:56 alpha ease-out", false, (12*3600 + 34*60 + 56) * 1000, t.TransitionEaseOut},

		// Transition ease-in
		{"00:02:00 alpha ease-in", false, 120_000, t.TransitionEaseIn},
		{"00:05:30 alpha ease-in", false, (5*60 + 30) * 1000, t.TransitionEaseIn},

		// Transition smooth
		{"00:03:00 alpha smooth", false, 180_000, t.TransitionSmooth},
		{"01:00:00 alpha smooth", false, 3_600_000, t.TransitionSmooth},

		// Empty transition (steady default)
		{"00:00:00 alpha", false, 0, t.TransitionSteady},
		{"00:10:00 alpha", false, 600_000, t.TransitionSteady},

		// Invalid transition types
		{"00:00:05 alpha invalid-transition", true, 0, t.TransitionSteady},
		{"00:00:05 alpha linear", true, 0, t.TransitionSteady},

		// Extra tokens after valid transition
		{"00:00:05 alpha steady extra", true, 0, t.TransitionSteady},
		{"00:00:05 alpha ease-in extra-token", true, 0, t.TransitionSteady},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		per, err := ctx.ParseTimeline(&presets)
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got none", test.line)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}
		if per == nil {
			ts.Errorf("For line '%s', expected non-nil period", test.line)
			continue
		}
		if per.Time != test.expectedMs {
			ts.Errorf("For line '%s', expected time %d but got %d", test.line, test.expectedMs, per.Time)
		}
		if per.Transition != test.expectedTransition {
			ts.Errorf("For line '%s', expected transition %v but got %v", test.line, test.expectedTransition, per.Transition)
		}
	}
}

func TestParseTimeline_TemplatePresetNotAllowed(ts *testing.T) {
	var presets []t.Preset

	// Create a template preset
	templatePreset, err := t.NewPreset("base-template", true, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating template preset: %v", err)
	}
	presets = append(presets, *templatePreset)

	// Create a normal preset for comparison
	normalPreset, err := t.NewPreset("alpha", false, nil)
	if err != nil {
		ts.Fatalf("unexpected error creating normal preset: %v", err)
	}
	presets = append(presets, *normalPreset)

	tests := []struct {
		name        string
		line        string
		expectError bool
	}{
		{
			name:        "template preset in timeline should fail",
			line:        "00:00:00 base-template",
			expectError: true,
		},
		{
			name:        "template preset with transition should fail",
			line:        "00:01:00 base-template steady",
			expectError: true,
		},
		{
			name:        "template preset with ease-out should fail",
			line:        "00:02:00 base-template ease-out",
			expectError: true,
		},
		{
			name:        "normal preset should succeed",
			line:        "00:00:00 alpha",
			expectError: false,
		},
		{
			name:        "normal preset with transition should succeed",
			line:        "00:01:00 alpha smooth",
			expectError: false,
		},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		per, err := ctx.ParseTimeline(&presets)

		if test.expectError {
			if err == nil {
				ts.Errorf("%s: expected error but got none for line '%s'", test.name, test.line)
				continue
			}
		} else {
			if err != nil {
				ts.Errorf("%s: unexpected error for line '%s': %v", test.name, test.line, err)
				continue
			}
			if per == nil {
				ts.Errorf("%s: expected non-nil period for line '%s'", test.name, test.line)
			}
		}
	}
}
