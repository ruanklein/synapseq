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

func TestHasComment(t *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"# This is a comment", true},
		{"// Another comment", false},
		{"No comment here", false},
		{"#Comment without space", true},
		{"   # Indented comment", true},
		{"## Double Comment!", true},
		{"  ## Indented double Comment!", true},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.HasComment()
		if result != test.expected {
			t.Errorf("For line '%s', expected HasComment() to be %v but got %v", test.line, test.expected, result)
		}
	}
}

func TestParseComment(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"# This is a comment", ""},
		{"// Another comment", ""},
		{"No comment here", ""},
		{"#Comment without space", ""},
		{"   # Indented comment", ""},
		{"## Double Comment!", "Double Comment!"},
		{"  ## Indented double Comment!", "Indented double Comment!"},
		{"##", " "},
		{"# First part // not a comment", ""},
	}

	for _, test := range tests {
		ctx := NewTextParser(test.line)
		result := ctx.ParseComment()
		if result != test.expected {
			t.Errorf("For line '%s', expected ParseComment() to be '%s' but got '%s'", test.line, test.expected, result)
		}
	}
}
