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

// NewFFmpegUnsafe creates an FFmpeg instance without validating the path.
// Useful for documentation examples and testing environments.
func NewFFmpegUnsafe(path string) *FFmpeg {
	if path == "" {
		path = "ffmpeg"
	}
	return &FFmpeg{baseUtility: baseUtility{path: path}}
}

// MP3 encodes streaming PCM into an MP3 file using ffmpeg.
func (fm *FFmpeg) MP3(appCtx *synapseq.AppContext, options *MP3Options) error {
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

	optsLine := [2]string{"-q:a", "0"} // Default to highest VBR quality (V0)
	if options != nil && options.Mode == MP3ModeCBR {
		optsLine[0] = "-b:a"
		optsLine[1] = "320k" // CBR at 320 kbps
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
		optsLine[0], optsLine[1],
		"-vn",
		outputFile,
	)

	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}

// OGG encodes streaming PCM into an OGG file using ffmpeg.
func (fm *FFmpeg) OGG(appCtx *synapseq.AppContext) error {
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

	// Vorbis quality scale is typically 0..10
	ffmpeg := fm.Command(
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
		"-c:a", "libvorbis",
		"-q:a", "10", // Highest quality
		"-vn",
		outputFile,
	)

	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}

// OPUS encodes streaming PCM into an OPUS file using ffmpeg.
func (fm *FFmpeg) OPUS(appCtx *synapseq.AppContext) error {
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

	// Use libopus with specified target bitrate
	ffmpeg := fm.Command(
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
		"-c:a", "libopus",
		"-b:a", "96k",
		"-vn",
		outputFile,
	)

	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}
