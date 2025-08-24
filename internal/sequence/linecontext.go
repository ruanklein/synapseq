package sequence

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type LineContext struct {
	Line   string   // Current line for processing
	Tokens []string // Tokens extracted from the line
	tkIdx  int      // Current token index
}

// Peek retrieves the next token without advancing the index
func (ctx *LineContext) Peek() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		return ctx.Tokens[ctx.tkIdx], true
	}
	return "", false
}

// NextToken retrieves the next token from the context
func (ctx *LineContext) NextToken() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		token := ctx.Tokens[ctx.tkIdx]
		ctx.tkIdx++
		return token, true
	}
	return "", false
}

// RewindToken moves the token index back by n positions
func (ctx *LineContext) RewindToken(n int) {
	if n <= 0 {
		return
	}
	if n > ctx.tkIdx {
		ctx.tkIdx = 0
		return
	}
	ctx.tkIdx -= n
}

// NextExpectOneOf checks if the next token is one of the expected values
func (ctx *LineContext) NextExpectOneOf(wants ...string) (string, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return "", fmt.Errorf("expected one of %v, got EOF: %s", wants, ctx.Line)
	}

	if slices.Contains(wants, tok) {
		return tok, nil
	}
	return "", fmt.Errorf("expected one of %v, got %q: %s", wants, tok, ctx.Line)
}

// NextFloat64Strict retrieves the next token as a float64, enforcing strict parsing
func (ctx *LineContext) NextFloat64Strict() (float64, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, fmt.Errorf("expected number, got EOF: %s", ctx.Line)
	}

	f, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %q", tok)
	}
	return f, nil
}

// NewLineContext creates a new LineContext for the given line
func NewLineContext(line string) *LineContext {
	return &LineContext{
		Line:   line,
		Tokens: strings.Fields(line),
		tkIdx:  0,
	}
}
