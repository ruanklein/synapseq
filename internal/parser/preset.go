package parser

import (
	"fmt"
	"strings"

	"github.com/ruanklein/synapseq/internal/types"
)

// isLetter checks if a byte is a letter
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// isDigit checks if a byte is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (ctx *ParserContext) IsPresetLine() bool {
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
func (ctx *ParserContext) ParsePresetLine() (*types.Preset, error) {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return nil, fmt.Errorf("expected preset name, got EOF: %s", ln)
	}

	presetName := strings.ToLower(tok)
	if presetName == types.BuiltinSilence {
		return nil, fmt.Errorf("cannot load %q built-in preset: %s", presetName, ln)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return nil, fmt.Errorf("unexpected token after definition: %q", unknown)
	}

	preset := &types.Preset{Name: presetName}
	preset.InitVoices()
	return preset, nil
}
