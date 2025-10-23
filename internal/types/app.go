/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

// FileFormat represents the format of the input/output file
type FileFormat int

const (
	FormatText FileFormat = iota
	FormatJSON
	FormatXML
	FormatYAML
)

// String returns the string representation of the FileFormat
func (ff FileFormat) String() string {
	switch ff {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	case FormatXML:
		return "xml"
	case FormatYAML:
		return "yaml"
	default:
		return "unknown"
	}
}

// AppMode represents the mode of the application
type AppMode int

const (
	ModeGenerate AppMode = iota
	ModeExtract
	ModeConvert
)

// String returns the string representation of the AppMode
func (am AppMode) String() string {
	switch am {
	case ModeGenerate:
		return "generate"
	case ModeExtract:
		return "extract"
	case ModeConvert:
		return "convert"
	default:
		return "unknown"
	}
}

// AppContext holds the application context
type AppContext struct {
	// Application mode
	Mode AppMode
	// Input file path
	InputFile string
	// Output file path
	OutputFile string
	// File format
	Format FileFormat
	// Quiet mode, suppress non-error output
	Quiet bool
	// Debug mode, no wav output
	Debug bool
	// Do not embed metadata in output WAV
	NoEmbedMetadata bool
	// Loaded sequence
	Sequence *Sequence
}
