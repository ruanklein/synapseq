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
		expectError  bool
	}{
		// Version flag
		{
			args:         []string{"cmd", "-version"},
			expected:     &CLIOptions{ShowVersion: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Help flag
		{
			args:         []string{"cmd", "-help"},
			expected:     &CLIOptions{ShowHelp: true},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Quiet flag
		{
			args:         []string{"cmd", "-quiet", "input.spsq", "output.wav"},
			expected:     &CLIOptions{Quiet: true},
			expectedArgs: []string{"input.spsq", "output.wav"},
			expectError:  false,
		},
		// Test flag
		{
			args:         []string{"cmd", "-test", "input.spsq"},
			expected:     &CLIOptions{Test: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// JSON format flag
		{
			args:         []string{"cmd", "-json", "input.json", "output.wav"},
			expected:     &CLIOptions{FormatJSON: true},
			expectedArgs: []string{"input.json", "output.wav"},
			expectError:  false,
		},
		// XML format flag
		{
			args:         []string{"cmd", "-xml", "input.xml", "output.wav"},
			expected:     &CLIOptions{FormatXML: true},
			expectedArgs: []string{"input.xml", "output.wav"},
			expectError:  false,
		},
		// YAML format flag
		{
			args:         []string{"cmd", "-yaml", "input.yaml", "output.wav"},
			expected:     &CLIOptions{FormatYAML: true},
			expectedArgs: []string{"input.yaml", "output.wav"},
			expectError:  false,
		},
		// Multiple flags combined
		{
			args:         []string{"cmd", "-quiet", "-json", "input.json", "output.wav"},
			expected:     &CLIOptions{Quiet: true, FormatJSON: true},
			expectedArgs: []string{"input.json", "output.wav"},
			expectError:  false,
		},
		// Test and quiet combined
		{
			args:         []string{"cmd", "-test", "-quiet", "input.spsq"},
			expected:     &CLIOptions{Test: true, Quiet: true},
			expectedArgs: []string{"input.spsq"},
			expectError:  false,
		},
		// XML format with quiet
		{
			args:         []string{"cmd", "-xml", "-quiet", "input.xml", "output.wav"},
			expected:     &CLIOptions{FormatXML: true, Quiet: true},
			expectedArgs: []string{"input.xml", "output.wav"},
			expectError:  false,
		},
		// No flags, just arguments
		{
			args:         []string{"cmd", "input.spsq", "output.wav"},
			expected:     &CLIOptions{},
			expectedArgs: []string{"input.spsq", "output.wav"},
			expectError:  false,
		},
		// No arguments at all
		{
			args:         []string{"cmd"},
			expected:     &CLIOptions{},
			expectedArgs: []string{},
			expectError:  false,
		},
		// Multiple format flags (should work, last one wins in flag package)
		{
			args:         []string{"cmd", "-json", "-xml", "input.xml", "output.wav"},
			expected:     &CLIOptions{FormatJSON: true, FormatXML: true},
			expectedArgs: []string{"input.xml", "output.wav"},
			expectError:  false,
		},
		// All boolean flags enabled
		{
			args:         []string{"cmd", "-quiet", "-test", "-json", "input.json"},
			expected:     &CLIOptions{Quiet: true, Test: true, FormatJSON: true},
			expectedArgs: []string{"input.json"},
			expectError:  false,
		},
		// Invalid flag should return error
		{
			args:         []string{"cmd", "-invalid"},
			expected:     nil,
			expectedArgs: nil,
			expectError:  true,
		},
		// Unknown flag with valid flags
		{
			args:         []string{"cmd", "-quiet", "-unknown", "input.spsq"},
			expected:     nil,
			expectedArgs: nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		os.Args = test.args
		opts, args, err := ParseFlags()

		if test.expectError {
			if err == nil {
				ts.Errorf("For args %v, expected error but got none", test.args)
			}
			continue
		}

		if err != nil {
			ts.Errorf("For args %v, unexpected error: %v", test.args, err)
			continue
		}

		if opts.ShowVersion != test.expected.ShowVersion {
			ts.Errorf("For args %v, ShowVersion: expected %v but got %v", test.args, test.expected.ShowVersion, opts.ShowVersion)
		}
		if opts.ShowHelp != test.expected.ShowHelp {
			ts.Errorf("For args %v, ShowHelp: expected %v but got %v", test.args, test.expected.ShowHelp, opts.ShowHelp)
		}
		if opts.Quiet != test.expected.Quiet {
			ts.Errorf("For args %v, Quiet: expected %v but got %v", test.args, test.expected.Quiet, opts.Quiet)
		}
		if opts.Test != test.expected.Test {
			ts.Errorf("For args %v, Test: expected %v but got %v", test.args, test.expected.Test, opts.Test)
		}
		if opts.FormatJSON != test.expected.FormatJSON {
			ts.Errorf("For args %v, FormatJSON: expected %v but got %v", test.args, test.expected.FormatJSON, opts.FormatJSON)
		}
		if opts.FormatXML != test.expected.FormatXML {
			ts.Errorf("For args %v, FormatXML: expected %v but got %v", test.args, test.expected.FormatXML, opts.FormatXML)
		}
		if opts.FormatYAML != test.expected.FormatYAML {
			ts.Errorf("For args %v, FormatYAML: expected %v but got %v", test.args, test.expected.FormatYAML, opts.FormatYAML)
		}

		if len(args) != len(test.expectedArgs) {
			ts.Errorf("For args %v, expected args %v but got %v", test.args, test.expectedArgs, args)
		} else {
			for i := range args {
				if args[i] != test.expectedArgs[i] {
					ts.Errorf("For args %v, expected args[%d] = %q but got %q", test.args, i, test.expectedArgs[i], args[i])
					break
				}
			}
		}
	}
}

func TestParseFlagsEdgeCases(ts *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test with stdin input
	os.Args = []string{"cmd", "-quiet", "-", "output.wav"}
	opts, args, err := ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for stdin input: %v", err)
	}
	if !opts.Quiet {
		ts.Errorf("expected Quiet=true for stdin test")
	}
	if len(args) != 2 || args[0] != "-" || args[1] != "output.wav" {
		ts.Errorf("expected args [\"-\", \"output.wav\"], got %v", args)
	}

	// Test with stdout output
	os.Args = []string{"cmd", "-json", "input.json", "-"}
	opts, args, err = ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for stdout output: %v", err)
	}
	if !opts.FormatJSON {
		ts.Errorf("expected FormatJSON=true for stdout test")
	}
	if len(args) != 2 || args[0] != "input.json" || args[1] != "-" {
		ts.Errorf("expected args [\"input.json\", \"-\"], got %v", args)
	}

	// Test with URL input
	os.Args = []string{"cmd", "-yaml", "https://example.com/sequence.yaml", "output.wav"}
	opts, args, err = ParseFlags()
	if err != nil {
		ts.Errorf("unexpected error for URL input: %v", err)
	}
	if !opts.FormatYAML {
		ts.Errorf("expected FormatYAML=true for URL test")
	}
	if len(args) != 2 || args[0] != "https://example.com/sequence.yaml" || args[1] != "output.wav" {
		ts.Errorf("expected URL args, got %v", args)
	}
}
