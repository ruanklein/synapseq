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
func (fm *FFmpeg) metadataArgs(metadata *info.Metadata) map[string]string {
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

// Convert encodes streaming PCM into the specified format using ffmpeg.
func (fm *FFmpeg) Convert(appCtx *synapseq.AppContext, format string, options *CodecOptions) error {
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

	sampleRate := appCtx.SampleRate()
	if format == "opus" && sampleRate != 48000 {
		return fmt.Errorf("opus format requires a sample rate of 48000 Hz")
	}

	// Base ffmpeg arguments
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(sampleRate),
		"-i", "pipe:0",
	}

	// Determine format and corresponding options
	switch format {
	case "mp3":
		args = append(args, []string{
			"-c:a", "libmp3lame",
		}...)

		// Determine MP3 encoding mode
		if options != nil && options.MP3Options != nil && options.MP3Options.Mode == MP3ModeCBR {
			args = append(args, []string{
				"-b:a", "320k", // CBR at 320 kbps
			}...)
		} else {
			args = append(args, []string{
				"-q:a", "0", // Highest VBR quality (V0)
			}...)
		}

		args = append(args, []string{
			"-f", "mp3",
		}...)
	case "ogg":
		args = append(args, []string{
			"-c:a", "libvorbis",
			"-q:a", "10", // Highest quality
			"-f", "ogg",
		}...)
	case "opus":
		args = append(args, []string{
			"-c:a", "libopus",
			"-b:a", "96k",
			"-f", "opus",
		}...)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	// Metadata embedding
	if !appCtx.UnsafeNoMetadata() && appCtx.Format() == "text" {
		rawContent := appCtx.RawContent()
		if rawContent == nil {
			return fmt.Errorf("raw content is nil for metadata embedding")
		}

		metadata, err := info.NewMetadata(rawContent)
		if err != nil {
			return fmt.Errorf("failed to create metadata: %v", err)
		}

		metaArgs := fm.metadataArgs(metadata)
		for key, value := range metaArgs {
			args = append(args, "-metadata", fmt.Sprintf("%s=%s", key, value))
		}
	}

	args = append(args, []string{
		"-vn",
		outputFile,
	}...)

	ffmpeg := fm.Command(args...)
	if err := startPipeCmd(ffmpeg, appCtx); err != nil {
		return err
	}

	return nil
}
