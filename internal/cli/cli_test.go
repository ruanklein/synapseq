/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package cli

import (
	"os"
	"testing"
)

func TestParseFlags(ts *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		args         []string
		expected     *CLIOptions
		expectedArgs []string
	}{
		{
			args:         []string{"cmd", "-version"},
			expected:     &CLIOptions{ShowVersion: true},
			expectedArgs: []string{},
		},
		{
			args:         []string{"cmd", "-quiet", "input.spsq", "output.wav"},
			expected:     &CLIOptions{Quiet: true},
			expectedArgs: []string{"input.spsq", "output.wav"},
		},
		{
			args:         []string{"cmd", "input.spsq", "output.wav"},
			expected:     &CLIOptions{},
			expectedArgs: []string{"input.spsq", "output.wav"},
		},
		{
			args:         []string{"cmd"},
			expected:     &CLIOptions{},
			expectedArgs: []string{},
		},
	}

	for _, test := range tests {
		os.Args = test.args
		opts, args, err := ParseFlags()
		if err != nil {
			ts.Errorf("For args %v, unexpected error: %v", test.args, err)
			continue
		}
		if *opts != *test.expected {
			ts.Errorf("For args %v, expected %+v but got %+v", test.args, test.expected, opts)
		}
		if len(args) != len(test.expectedArgs) {
			ts.Errorf("For args %v, expected args %v but got %v", test.args, test.expectedArgs, args)
		} else {
			for i := range args {
				if args[i] != test.expectedArgs[i] {
					ts.Errorf("For args %v, expected args %v but got %v", test.args, test.expectedArgs, args)
					break
				}
			}
		}
	}

	// Test invalid flag parsing should return an error (no panic)
	os.Args = []string{"cmd", "-invalid"}
	if _, _, err := ParseFlags(); err == nil {
		ts.Errorf("expected error for invalid flag, got nil")
	}
}
