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
