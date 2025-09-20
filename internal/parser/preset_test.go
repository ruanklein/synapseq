/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"testing"
)

func TestHasPreset(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"MyPreset", true},
		{" AnotherPreset", false},
		{"preset1", true},
		{"  preset2", false},
		{"123Preset", false},
		{"", false},
		{"   ", false},
		{"%Preset", false},
		{"Preset_", true},
		{"preset-01", true},
		{"preset-", true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasPreset()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasPreset() to return %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParsePreset(ts *testing.T) {
	tests := []struct {
		line          string
		expectedName  string
		expectedError bool
	}{
		{"MyPreset", "mypreset", false},
		{"%AnotherPreset", "", true},
		{"preset1", "preset1", false},
		{"123Preset", "", true},
		{"", "", true},
		{"   ", "", true},
		{"Preset_", "preset_", false},
		{"preset-01", "preset-01", false},
		{"preset-", "preset-", false},
		{"silence", "", true}, // reserved name
		{"Pre$et", "", true},  // invalid character
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		preset, err := ctx.ParsePreset()
		if test.expectedError {
			if err == nil {
				ts.Errorf("For line '%s', expected an error but got none", test.line)
			}
			continue
		}
		if err != nil {
			ts.Errorf("For line '%s', did not expect an error but got: %v", test.line, err)
			continue
		}
		if preset.String() != test.expectedName {
			ts.Errorf("For line '%s', expected preset name '%s' but got '%s'", test.line, test.expectedName, preset.String())
		}
	}
}
