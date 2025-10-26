/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	"fmt"
	"io"

	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// AppContext holds the configuration for the application.
// It provides a safe, immutable context for sequence processing.
// Methods that modify the context return a new instance.
type AppContext struct {
	inputFile        string
	outputFile       string
	format           t.FileFormat
	unsafeNoMetadata bool
	statusOutput     io.Writer
	sequence         *t.Sequence
}

// NewAppContext creates a new AppContext instance.
//
// Parameters:
//   - inputFile: path to the input sequence file (can be local path, stdin "-", or HTTP/HTTPS URL)
//   - outputFile: path to the output WAV file (local path only)
//   - format: file format, one of: "text", "json", "xml", "yaml"
//
// Returns an error if the format is invalid.
func NewAppContext(inputFile, outputFile, format string) (*AppContext, error) {
	var fileFormat t.FileFormat
	switch format {
	case "text":
		fileFormat = t.FormatText
	case "json":
		fileFormat = t.FormatJSON
	case "xml":
		fileFormat = t.FormatXML
	case "yaml":
		fileFormat = t.FormatYAML
	default:
		return nil, fmt.Errorf("invalid file format: %s", format)
	}

	return &AppContext{
		inputFile:        inputFile,
		outputFile:       outputFile,
		format:           fileFormat,
		unsafeNoMetadata: false,
		statusOutput:     nil,
	}, nil
}

// InputFile returns the input file path.
func (ac *AppContext) InputFile() string {
	return ac.inputFile
}

// OutputFile returns the output file path.
func (ac *AppContext) OutputFile() string {
	return ac.outputFile
}

// Format returns the file format as a string.
func (ac *AppContext) Format() string {
	return ac.format.String()
}

// Verbose returns whether verbose mode is enabled.
// When true, status output will be written to the configured writer.
func (ac *AppContext) Verbose() bool {
	return ac.statusOutput != nil
}

// UnsafeNoMetadata returns whether the unsafe no metadata flag is set.
// When true, metadata validation is disabled for text format files.
func (ac *AppContext) UnsafeNoMetadata() bool {
	return ac.unsafeNoMetadata
}

// WithVerbose returns a new AppContext with verbose mode enabled.
// Status output will be written to the provided writer (typically os.Stderr).
//
// Example:
//
//	ctx = ctx.WithVerbose(os.Stderr)
func (ac *AppContext) WithVerbose(data io.Writer) *AppContext {
	newCtx := *ac
	newCtx.statusOutput = data
	return &newCtx
}

// WithUnsafeNoMetadata returns a new AppContext with the unsafe no metadata flag set.
// This option is only available for text format files.
//
// WARNING: This option disables metadata validation and should only be used
// when you understand the implications.
//
// Returns an error if called on non-text format.
func (ac *AppContext) WithUnsafeNoMetadata() (*AppContext, error) {
	if ac.format != t.FormatText {
		return nil, fmt.Errorf("unsafe no metadata can only be set for text format")
	}

	newCtx := *ac
	newCtx.unsafeNoMetadata = true
	return &newCtx, nil
}
