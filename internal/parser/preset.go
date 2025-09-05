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

// HasPreset checks if the current line is a preset definition
func (ctx *TextParser) HasPreset() bool {
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	// first char of line must be a letter
	first := ctx.Line.Raw[0]
	if !isLetter(first) {
		return false
	}

	// can contain letters, digits, '_' or '-'
	for i := 1; i < len(tok); i++ {
		ch := tok[i]
		if !(isLetter(ch) || isDigit(ch) || ch == '_' || ch == '-') {
			return false
		}
	}

	return true
}

// ParsePreset extracts and returns a Preset from the current line context
func (ctx *TextParser) ParsePreset() (*t.Preset, error) {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected preset name, got EOF: %s", ln)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after preset definition: %q", unknown)
	}

	preset := &t.Preset{Name: strings.ToLower(tok)}
	preset.InitVoices(t.VoiceOff)
	return preset, nil
}
