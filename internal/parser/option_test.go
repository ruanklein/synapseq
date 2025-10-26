/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

func TestHasOption(ts *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{fmt.Sprintf("%svolume 50", t.KeywordOption), true},
		{fmt.Sprintf("%ssamplerate 48000", t.KeywordOption), true},
		{fmt.Sprintf("   %sgainlevel medium", t.KeywordOption), false},
		{fmt.Sprintf("background file.wav %s", t.KeywordComment), false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasOption()
		if result != test.expected {
			ts.Errorf("For line '%s', expected HasOption() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseOption(ts *testing.T) {
	// Fake path for testing background option
	backgroundFile := "noise.wav"
	cwd, err := os.Getwd()
	if err != nil {
		ts.Errorf("cannot get current working directory")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		ts.Errorf("cannot get user home directory")
	}
	// Test valid options
	tests := []struct {
		line     string
		expected t.SequenceOptions
	}{
		{fmt.Sprintf("%svolume 50", t.KeywordOption), t.SequenceOptions{Volume: 50}},
		{fmt.Sprintf("%ssamplerate 48000", t.KeywordOption), t.SequenceOptions{SampleRate: 48000}},
		{fmt.Sprintf("%sgainlevel low", t.KeywordOption), t.SequenceOptions{GainLevel: t.GainLevelLow}},
		{fmt.Sprintf("%sbackground testdata/%s", t.KeywordOption, backgroundFile), t.SequenceOptions{BackgroundPath: filepath.Join(cwd+"/testdata/", backgroundFile)}},
		{fmt.Sprintf("%sbackground ~/Downloads/%s", t.KeywordOption, backgroundFile), t.SequenceOptions{BackgroundPath: filepath.Join(homeDir+"/Downloads/", backgroundFile)}},
	}

	for _, test := range tests {
		option := t.SequenceOptions{}
		ctx := NewTextParser(test.line)
		if err := ctx.ParseOption(&option); err != nil {
			ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			continue
		}

		if option != test.expected {
			ts.Errorf("For line '%s', expected option %+v but got %+v", test.line, test.expected, option)
		}
	}
}
