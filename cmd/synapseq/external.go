/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"

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
func externalMp3(ffmpegPath, mode string, appCtx *synapseq.AppContext) error {
	ffmpeg, err := external.NewFFmpeg(ffmpegPath)
	if err != nil {
		return err
	}

	var mp3Mode external.MP3Mode
	switch mode {
	case "vbr":
		mp3Mode = external.MP3ModeVBR
	case "cbr":
		mp3Mode = external.MP3ModeCBR
	default:
		return fmt.Errorf("invalid MP3 mode: %s", mode)
	}

	if err := ffmpeg.MP3(appCtx, &external.MP3Options{Mode: mp3Mode}); err != nil {
		return err
	}

	return nil
}

// externalOgg encodes streaming PCM into an MP3 file using external utility
func externalOgg(ffmpegPath string, appCtx *synapseq.AppContext) error {
	ffmpeg, err := external.NewFFmpeg(ffmpegPath)
	if err != nil {
		return err
	}

	if err := ffmpeg.OGG(appCtx); err != nil {
		return err
	}

	return nil
}
