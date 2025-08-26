package parser

import (
	"fmt"
	"strings"

	t "github.com/ruanklein/synapseq/internal/types"
)

// isLetter checks if a byte is a letter
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// isDigit checks if a byte is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// isPreset checks if a string is a valid preset name
func isPreset(s string) bool {
	if len(s) == 0 {
		return false
	}

	if !isLetter(s[0]) {
		return false
	}

	for i := 1; i < len(s); i++ {
		ch := s[i]
		if !(isLetter(ch) || isDigit(ch) || ch == '_' || ch == '-') {
			return false
		}
	}

	return true
}

// IsPresetLine checks if the current line is a preset definition
func (ctx *ParserContext) IsPresetLine() bool {
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	if ctx.Line.Raw[0] == ' ' {
		return false
	}

	return isPreset(tok)
}

// ParsePreset extracts and returns a Preset from the current line context
func (ctx *ParserContext) ParsePresetLine() (*t.Preset, error) {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected preset name, got EOF: %s", ln)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after definition: %q", unknown)
	}

	preset := &t.Preset{Name: strings.ToLower(tok)}
	preset.InitVoices()
	return preset, nil
}
