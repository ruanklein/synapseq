package parser

import (
	"fmt"
	"strings"

	t "github.com/ruanklein/synapseq/internal/types"
)

// HasComment checks if the first element is a comment
func (ctx *TextParser) HasComment() bool {
	tok, ok := ctx.Line.Peek()
	return ok && string(tok[0]) == t.KeywordComment
}

// ParseComment extracts and prints the comment from the elements
func (ctx *TextParser) ParseComment() string {
	tok, ok := ctx.Line.Peek()
	if !ok || string(tok[0]) != t.KeywordComment {
		return ""
	}
	if len(tok) >= 2 && string(tok[1]) == t.KeywordComment {
		return fmt.Sprintf("%s %s", tok[2:], strings.Join(ctx.Line.Tokens[1:], " "))
	}
	return ""
}
