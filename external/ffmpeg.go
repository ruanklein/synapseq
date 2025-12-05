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
	"github.com/ruanklein/synapseq/v3/internal/info"
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

// metadataArgs returns ffmpeg arguments for embedding metadata.
func metadataArgs(metadata *info.Metadata) map[string]string {
	if metadata == nil {
		return nil
	}

	return map[string]string{
		"synapseq_id":        metadata.ID(),
		"synapseq_generated": metadata.Generated(),
		"synapseq_version":   metadata.Version(),
		"synapseq_platform":  metadata.Platform(),
		"synapseq_content":   metadata.Content(),
	}
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

	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
		"-c:a", "libmp3lame",
		"-f", "mp3",
		optsLine[0], optsLine[1],
	}

	if !appCtx.UnsafeNoMetadata() && appCtx.Format() == "text" {
		rawContent := appCtx.RawContent()
		if rawContent == nil {
			return fmt.Errorf("raw content is nil for metadata embedding")
		}

		metadata, err := info.NewMetadata(rawContent)
		if err != nil {
			return fmt.Errorf("failed to create metadata: %v", err)
		}

		metaArgs := metadataArgs(metadata)
		for key, value := range metaArgs {
			args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
		}
	}

	args = append(args, []string{
		"-vn",
		outputFile,
	}...)

	// ffmpeg command for highest MP3 quality (LAME V0)
	ffmpeg := fm.Command(args...)
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
		"-f", "ogg",
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
		"-f", "opus",
		"-vn",
		outputFile,
	)

	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}
