//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/external"
)

// externalPlay invokes utility tool to play from streaming audio input
func externalPlay(ffplayPath string, appCtx *synapseq.AppContext) error {
	ffplay, err := external.NewFFPlay(ffplayPath)
	if err != nil {
		return err
	}

	if err := ffplay.Play(appCtx); err != nil {
		return err
	}

	return nil
}

// externalMp3 encodes streaming PCM into an MP3 file using external utility
func externalMp3(ffmpegPath string, appCtx *synapseq.AppContext) error {
	ffmpeg, err := external.NewFFmpeg(ffmpegPath)
	if err != nil {
		return err
	}

	if err := ffmpeg.Convert(appCtx, "mp3"); err != nil {
		return err
	}

	return nil
}

// externalExtractTextSequence extracts text sequence from input file using ffprobe
func externalExtractTextSequence(ffprobePath string, inputFile string) (string, error) {
	ffprobe, err := external.NewFFprobe(ffprobePath)
	if err != nil {
		return "", err
	}

	content, err := ffprobe.ExtractTextSequence(inputFile)
	if err != nil {
		return "", err
	}

	return content, nil
}

// externalSaveExtractedTextSequence saves extracted text sequence to output file using ffprobe
func externalSaveExtractedTextSequence(ffprobePath, inputFile, outputFile string) error {
	ffprobe, err := external.NewFFprobe(ffprobePath)
	if err != nil {
		return err
	}

	if err := ffprobe.SaveExtractedTextSequence(inputFile, outputFile); err != nil {
		return err
	}

	return nil
}
