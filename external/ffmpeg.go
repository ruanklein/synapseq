/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import (
	"fmt"
	"os"
	"strconv"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// FFmpeg represents the ffmpeg external tool
type FFmpeg struct{ baseUtility }

// NewFFmpeg creates a new FFmpeg instance with given ffmpeg path
func NewFFmpeg(ffmpegPath string) (*FFmpeg, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	util, err := newUtility(ffmpegPath)
	if err != nil {
		return nil, err
	}

	return &FFmpeg{baseUtility: *util}, nil
}

// MP3 encodes streaming PCM into an MP3 file using ffmpeg.
func (fm *FFmpeg) MP3(appCtx *synapseq.AppContext) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	// Remove existing output file if it exists
	outputFile := appCtx.OutputFile()
	if _, err := os.Stat(outputFile); err == nil {
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("failed to remove existing output file: %v", err)
		}
	}

	// ffmpeg command for highest MP3 quality (LAME V0)
	ffmpeg := fm.Command(
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
		"-c:a", "libmp3lame",
		"-q:a", "0", // Highest VBR quality (V0)
		"-vn",
		outputFile,
	)

	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}
