/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package external

import (
	"fmt"
	"os/exec"
	"strconv"

	synapseq "github.com/ruanklein/synapseq/v3/core"
)

// FFplay represents the ffplay external tool
type FFplay struct{ path string }

// NewFFPlay creates a new FFplay instance with given ffplay path
func NewFFPlay(ffplayPath string) (*FFplay, error) {
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	path, err := newUtility(ffplayPath)
	if err != nil {
		return nil, err
	}

	return &FFplay{
		path: path,
	}, nil
}

// Play invokes ffplay to play from streaming audio input
func (fp *FFplay) Play(appCtx *synapseq.AppContext) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	ffplay := exec.Command(
		fp.path,
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
