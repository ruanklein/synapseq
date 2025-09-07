package parser

import (
	"fmt"

	t "github.com/ruanklein/synapseq/internal/types"
)

// HasPreset checks if the current line is a preset definition
func (ctx *TextParser) HasPreset() bool {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	ch := tok[0]
	if ln[0] != ' ' && ((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
		return true
	}

	return false
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

	preset, err := t.NewPreset(tok)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return preset, nil
}
