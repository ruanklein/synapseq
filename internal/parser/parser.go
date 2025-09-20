/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package parser

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// TextParser holds the context for parsing
type TextParser struct {
	Line lineContext // Context for the current line
}

// lineContext holds the context for the current line being parsed
type lineContext struct {
	Raw    string   // Raw line text
	Tokens []string // Tokens extracted from the line
	tkIdx  int      // Current token index
}

// Peek retrieves the next token without advancing the index
func (ctx *lineContext) Peek() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		return ctx.Tokens[ctx.tkIdx], true
	}
	return "", false
}

// NextToken retrieves the next token from the context
func (ctx *lineContext) NextToken() (string, bool) {
	if ctx.tkIdx < len(ctx.Tokens) {
		token := ctx.Tokens[ctx.tkIdx]
		ctx.tkIdx++
		return token, true
	}
	return "", false
}

// RewindToken moves the token index back by n positions
func (ctx *lineContext) RewindToken(n int) {
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
func (ctx *lineContext) NextExpectOneOf(wants ...string) (string, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return "", fmt.Errorf("expected one of %v, got EOF: %s", wants, ctx.Raw)
	}

	if slices.Contains(wants, tok) {
		return tok, nil
	}
	return "", fmt.Errorf("expected one of %v, got %q: %s", wants, tok, ctx.Raw)
}

// NextFloat64Strict retrieves the next token as a float64, enforcing strict parsing
func (ctx *lineContext) NextFloat64Strict() (float64, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, fmt.Errorf("expected float, got EOF: %s", ctx.Raw)
	}

	f, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float: %q", tok)
	}

	// Reject NaN and Inf values
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid float (NaN or Inf): %q", tok)
	}

	// Reject scientific notation (e.g., 1e10)
	if strings.ContainsAny(tok, "eE") {
		return 0, fmt.Errorf("scientific notation not allowed: %q", tok)
	}

	return f, nil
}

// NextIntStrict retrieves the next token as an int, enforcing strict parsing
func (ctx *lineContext) NextIntStrict() (int, error) {
	tok, ok := ctx.NextToken()
	if !ok {
		return 0, fmt.Errorf("expected integer, got EOF: %s", ctx.Raw)
	}

	i, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %q", tok)
	}
	return i, nil
}

// NewTextParser creates a new TextParser for the given line
func NewTextParser(line string) *TextParser {
	return &TextParser{
		Line: lineContext{
			Raw:    line,
			Tokens: strings.Fields(line),
			tkIdx:  0,
		},
	}
}
