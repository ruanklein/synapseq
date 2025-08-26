package parser

import "fmt"

// parseTime parses a time string in HH:MM:SS format to milliseconds
func parseTime(s string) (int, error) {
	var hh, mm, ss int
	if _, err := fmt.Sscanf(s, "%d:%d:%d", &hh, &mm, &ss); err != nil {
		return 0, err
	}

	if ss < 0 || mm < 0 || hh < 0 || ss > 59 || mm > 59 || hh > 23 {
		return 0, fmt.Errorf("invalid time")
	}

	return (hh*3600 + mm*60 + ss) * 1000, nil
}

func (ctx *ParserContext) IsTimeline() bool {
	tok, ok := ctx.Line.Peek()
	if !ok {
		return false
	}

	if ctx.Line.Raw[0] == ' ' {
		return false
	}

	if _, err := parseTime(tok); err != nil {
		return false
	}

	return true
}

func (ctx *ParserContext) ParseTimeline() (string, error) {
	ln := ctx.Line.Raw
	tok, ok := ctx.Line.NextToken()
	if !ok {
		return "", fmt.Errorf("expected time, got EOF: %s", ln)
	}

	timeMs, err := parseTime(tok)
	if err != nil {
		return "", fmt.Errorf("invalid time %q: %s", tok, ln)
	}

	tok, ok = ctx.Line.NextToken()
	if !ok {
		return "", fmt.Errorf("expected preset name, got EOF: %s", ln)
	}

	if !isPreset(tok) {
		return "", fmt.Errorf("invalid preset name %q: %s", tok, ln)
	}

	unknown, ok := ctx.Line.Peek()
	if ok {
		return "", fmt.Errorf("unexpected token after definition: %q", unknown)
	}

	return fmt.Sprintf("Parsed timeline: time=%d, preset=%s\n", timeMs, tok), nil
}
