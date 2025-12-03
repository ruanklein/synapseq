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

	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/external"
)

// play invokes utility tool to play from streaming audio input
func play(playerPath, inputFile, format string, quiet bool) error {
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

	ffplay, err := external.NewFFPlay(playerPath)
	if err != nil {
		return err
	}

	if err := ffplay.Play(appCtx); err != nil {
		return err
	}

	return nil
}

// mp3 encodes streaming PCM into an MP3 file using external utility
func mp3(converterPath, inputFile, outputFile, format string, quiet bool) error {
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

	ffmpeg, err := external.NewFFmpeg(converterPath)
	if err != nil {
		return err
	}

	if err := ffmpeg.MP3(appCtx); err != nil {
		return err
	}

	return nil
}
