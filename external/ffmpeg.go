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
	"os/exec"
	"strconv"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// FFmpeg represents the ffmpeg external tool
type FFmpeg struct{ *externalTool }

// NewFFmpeg creates a new FFmpeg instance with given ffmpeg path
func NewFFmpeg(ffmpegPath string) (*FFmpeg, error) {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	et, err := newUtility(ffmpegPath)
	if err != nil {
		return nil, err
	}

	return &FFmpeg{
		externalTool: et,
	}, nil
}

// Path returns the path to the ffmpeg executable
func (fm *FFmpeg) Path() string {
	return fm.utilityPath
}

// MP3 encodes streaming PCM into an MP3 file using ffmpeg.
func (fm *FFmpeg) MP3(appCtx *synapseq.AppContext) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	// ffmpeg command for highest MP3 quality (LAME V0)
	ffmpeg := exec.Command(
		fm.utilityPath,
		"-hide_banner",
		"-loglevel", "error",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
		"-c:a", "libmp3lame",
		"-q:a", "0", // Highest VBR quality (V0)
		"-vn",
		appCtx.OutputFile(),
	)

	stdin, err := ffmpeg.StdinPipe()
	if err != nil {
		return err
	}

	ffmpeg.Stdout = os.Stdout
	ffmpeg.Stderr = os.Stderr

	if err := ffmpeg.Start(); err != nil {
		stdin.Close()
		return err
	}

	streamErr := appCtx.Stream(stdin)

	stdin.Close()
	waitErr := ffmpeg.Wait()

	if streamErr != nil {
		return streamErr
	}
	if waitErr != nil {
		return waitErr
	}

	return nil
}
