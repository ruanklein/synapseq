/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package types

const (
	// MaxTextFileSize is the maximum allowed size for text files (32KB)
	MaxTextFileSize = 32 * 1024
	// MaxBackgroundFileSize is the maximum allowed size for background files (10MB)
	MaxBackgroundFileSize = 10 * 1024 * 1024
	// MaxStructuredFileSize is the maximum allowed size for structured files (128KB)
	MaxStructuredFileSize = 128 * 1024
)

// FileFormat represents the format of the input/output file
type FileFormat int

const (
	FormatText FileFormat = iota
	FormatJSON
	FormatXML
	FormatYAML
	FormatWAV
	FormatUnknown
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
	case FormatWAV:
		return "wav"
	default:
		return "unknown"
	}
}
