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

	t "github.com/ruanklein/synapseq/internal/types"
)

// AppContext holds the configuration for the application
type AppContext struct {
	inputFile        string
	outputFile       string
	format           t.FileFormat
	unsafeNoMetadata bool
	statusOutput     io.Writer
	sequence         *t.Sequence
}

// NewAppContext creates a new AppContext instance
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

// InputFile returns the input file path
func (ac *AppContext) InputFile() string {
	return ac.inputFile
}

// OutputFile returns the output file path
func (ac *AppContext) OutputFile() string {
	return ac.outputFile
}

// Format returns the file format
func (ac *AppContext) Format() string {
	return ac.format.String()
}

// Verbose returns whether the application is in verbose mode
func (ac *AppContext) Verbose() bool {
	return ac.statusOutput != nil
}

// WithVerbose returns a new AppContext with verbose mode enabled
func (ac *AppContext) WithVerbose(data io.Writer) *AppContext {
	newCtx := *ac
	newCtx.statusOutput = data
	return &newCtx
}

// WithUnsafeNoMetadata returns a new AppContext with the unsafe no metadata flag set
func (ac *AppContext) WithUnsafeNoMetadata() (*AppContext, error) {
	if ac.format != t.FormatText {
		return nil, fmt.Errorf("unsafe no metadata can only be set for text format")
	}

	newCtx := *ac
	newCtx.unsafeNoMetadata = true
	return &newCtx, nil
}
