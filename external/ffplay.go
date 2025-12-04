/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import (
	"fmt"
	"strconv"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// FFplay represents the ffplay external tool
type FFplay struct{ baseUtility }

// NewFFPlay creates a new FFplay instance with given ffplay path
func NewFFPlay(ffplayPath string) (*FFplay, error) {
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	util, err := newUtility(ffplayPath)
	if err != nil {
		return nil, err
	}

	return &FFplay{baseUtility: *util}, nil
}

// NewFFPlayUnsafe creates an FFplay instance without validating the path.
// Useful for documentation examples and testing environments.
func NewFFPlayUnsafe(path string) *FFplay {
	if path == "" {
		path = "ffplay"
	}
	return &FFplay{baseUtility: baseUtility{path: path}}
}

// Play invokes ffplay to play from streaming audio input
func (fp *FFplay) Play(appCtx *synapseq.AppContext) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	ffplay := fp.Command(
		"-nodisp",
		"-hide_banner",
		"-loglevel", "error",
		"-autoexit",
		"-f", "s16le",
		"-ch_layout", "stereo",
		"-ar", strconv.Itoa(appCtx.SampleRate()),
		"-i", "pipe:0",
	)

	if err := startPipeCmd(ffplay, appCtx); err != nil {
		return err
	}

	return nil
}
