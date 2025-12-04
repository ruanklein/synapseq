/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

// MP3Mode represents the MP3 encoding mode
type MP3Mode int

const (
	// MP3ModeVBR represents Variable Bit Rate encoding with 0 as highest quality
	MP3ModeVBR MP3Mode = iota
	// MP3ModeCBR represents Constant Bit Rate encoding with 320 kbps
	MP3ModeCBR
)

// MP3Options holds options for MP3 encoding
type MP3Options struct {
	// Mode specifies the MP3 encoding mode (VBR or CBR)
	Mode MP3Mode
}
