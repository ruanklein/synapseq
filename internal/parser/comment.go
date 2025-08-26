package parser

import (
	"fmt"
	"strings"
)

// IsCommentLine checks if the first element is a comment
func (ctx *ParserContext) IsCommentLine() bool {
	tok, ok := ctx.Line.Peek()
	return ok && string(tok[0]) == keywordComment
}

// ParseComment extracts and prints the comment from the elements
func (ctx *ParserContext) ParseCommentLine() string {
	tok, ok := ctx.Line.Peek()
	if !ok || string(tok[0]) != keywordComment {
		return ""
	}
	if len(tok) >= 2 && string(tok[1]) == keywordComment {
		return fmt.Sprintf("%s %s", tok[2:], strings.Join(ctx.Line.Tokens[1:], " "))
	}
	return ""
}
