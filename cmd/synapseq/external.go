/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// externalTool holds paths to external Unix tools like ffmpeg and ffplay
type externalTool struct {
	ffmpegPath string
	ffplayPath string
}

// newExternalTool creates a new externalTool instance with given paths
func newExternalTool(ffmpegPath, ffplayPath string) *externalTool {
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	return &externalTool{
		ffmpegPath: ffmpegPath,
		ffplayPath: ffplayPath,
	}
}

// play invokes ffplay to play from streaming audio input
func (et *externalTool) play(inputFile, format string, quiet bool) error {
	appCtx, err := synapseq.NewAppContext(inputFile, "", format)
	if err != nil {
		return err
	}

	if err := appCtx.LoadSequence(); err != nil {
		return err
	}

	if !quiet {
		appCtx = appCtx.WithVerbose(os.Stderr)

		for _, c := range appCtx.Comments() {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	ffplay := exec.Command(
		et.ffplayPath,
		"-nodisp",
		"-hide_banner",
		"-loglevel", "error",
		"-autoexit",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
	)

	stdin, err := ffplay.StdinPipe()
	if err != nil {
		return err
	}

	ffplay.Stdout = os.Stdout
	ffplay.Stderr = os.Stderr

	if err := ffplay.Start(); err != nil {
		stdin.Close()
		return err
	}

	streamErr := appCtx.Stream(stdin)

	stdin.Close()

	waitErr := ffplay.Wait()

	if streamErr != nil {
		return streamErr
	}

	if waitErr != nil {
		return waitErr
	}

	return nil
}

// mp3 encodes streaming PCM into an MP3 file using ffmpeg.
func (et *externalTool) mp3(inputFile, outputFile, format string, quiet bool) error {
	appCtx, err := synapseq.NewAppContext(inputFile, outputFile, format)
	if err != nil {
		return err
	}

	if err := appCtx.LoadSequence(); err != nil {
		return err
	}

	if !quiet {
		appCtx = appCtx.WithVerbose(os.Stderr)

		for _, c := range appCtx.Comments() {
			fmt.Fprintf(os.Stderr, "> %s\n", c)
		}
	}

	// ffmpeg command for highest MP3 quality (LAME V0)
	ffmpeg := exec.Command(
		et.ffmpegPath,
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
