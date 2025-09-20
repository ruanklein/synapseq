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

func TestPeek(ts *testing.T) {
	tests := []struct {
		line          string
		expectedToken string
		expectedOk    bool
	}{
		{"waveform sine tone 440 binaural 10 amplitude 04", "waveform", true},
		{"   waveform sine tone 440 binaural 10 amplitude 04", "waveform", true},
		{"", "", false},
		{"   ", "", false},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		token, ok := ctx.Line.Peek()
		if token != test.expectedToken || ok != test.expectedOk {
			ts.Errorf("For line '%s', expected Peek() to return ('%s', %v) but got ('%s', %v)", test.line, test.expectedToken, test.expectedOk, token, ok)
		}
	}
}

func TestNextToken(ts *testing.T) {
	tests := []struct {
		line           string
		expectedTokens []string
	}{
		{"waveform sine tone 440 binaural 10 amplitude 04", []string{"waveform", "sine", "tone", "440", "binaural", "10", "amplitude", "04"}},
		{"   waveform sine tone 440 binaural 10 amplitude 04", []string{"waveform", "sine", "tone", "440", "binaural", "10", "amplitude", "04"}},
		{"", []string{}},
		{"   ", []string{}},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		var tokens []string
		for {
			token, ok := ctx.Line.NextToken()
			if !ok {
				break
			}
			tokens = append(tokens, token)
		}
		if len(tokens) != len(test.expectedTokens) {
			ts.Errorf("For line '%s', expected %d tokens but got %d", test.line, len(test.expectedTokens), len(tokens))
			continue
		}
		for i, expectedToken := range test.expectedTokens {
			if tokens[i] != expectedToken {
				ts.Errorf("For line '%s', expected token %d to be '%s' but got '%s'", test.line, i, expectedToken, tokens[i])
			}
		}
	}
}

func TestRewindToken(ts *testing.T) {
	line := "waveform sine tone 440 binaural 10 amplitude 04"
	ctx := NewTextParser(line)

	// Read first three tokens
	for range 3 {
		_, ok := ctx.Line.NextToken()
		if !ok {
			ts.Errorf("Unexpected EOF while reading tokens")
		}
	}

	// Rewind two tokens
	ctx.Line.RewindToken(2)

	// Next token should be the second token
	token, ok := ctx.Line.NextToken()
	if !ok || token != "sine" {
		ts.Errorf("Expected 'sine' after rewind, got '%s'", token)
	}

	// Rewind more than available tokens
	ctx.Line.RewindToken(10)

	// Next token should be the first token
	token, ok = ctx.Line.NextToken()
	if !ok || token != "waveform" {
		ts.Errorf("Expected 'waveform' after rewind to start, got '%s'", token)
	}
}

func TestNextExpectOneOf(ts *testing.T) {
	tests := []struct {
		line          string
		wants         []string
		expectedToken string
		expectError   bool
	}{
		{"waveform sine tone 440 binaural 10 amplitude 04", []string{"waveform", "noise"}, "waveform", false},
		{"noise pink amplitude 40", []string{"background", "noise"}, "noise", false},
		{"noise white amplitude 05", []string{"triangle", "background"}, "", true},
		{"background amplitude 50", []string{"noise", "waveform"}, "", true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		token, err := ctx.Line.NextExpectOneOf(test.wants...)
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got token '%s'", test.line, token)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if token != test.expectedToken {
				ts.Errorf("For line '%s', expected token '%s' but got '%s'", test.line, test.expectedToken, token)
			}
		}
	}
}

func TestFloat64Strict(ts *testing.T) {
	tests := []struct {
		line          string
		expectedValue float64
		expectError   bool
	}{
		{"440.5", 440.5, false},
		{"   123.456   ", 123.456, false},
		{"notanumber", 0, true},
		{"123abc", 0, true},
		{"", 0, true},
		{"   ", 0, true},
		{"NaN", 0, true},
		{"Inf", 0, true},
		{"-Inf", 0, true},
		{"1e10", 1e10, true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		value, err := ctx.Line.NextFloat64Strict()
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got value %f", test.line, value)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if value != test.expectedValue {
				ts.Errorf("For line '%s', expected value %f but got %f", test.line, test.expectedValue, value)
			}
		}
	}
}

func TestNextIntStrict(ts *testing.T) {
	tests := []struct {
		line          string
		expectedValue int
		expectError   bool
	}{
		{"48000", 48000, false},
		{"80", 80, false},
		{"   123   ", 123, false},
		{"notanumber", 0, true},
		{"123abc", 0, true},
		{"", 0, true},
		{"   ", 0, true},
		{"12.34", 0, true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		value, err := ctx.Line.NextIntStrict()
		if test.expectError {
			if err == nil {
				ts.Errorf("For line '%s', expected error but got value %d", test.line, value)
			}
		} else {
			if err != nil {
				ts.Errorf("For line '%s', unexpected error: %v", test.line, err)
			} else if value != test.expectedValue {
				ts.Errorf("For line '%s', expected value %d but got %d", test.line, test.expectedValue, value)
			}
		}
	}
}
