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
	alpha, err := t.NewPreset("alpha")
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
