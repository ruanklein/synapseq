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

// FFplay represents the ffplay external tool
type FFplay struct{ *externalTool }

// NewFFPlay creates a new FFplay instance with given ffplay path
func NewFFPlay(ffplayPath string) (*FFplay, error) {
	if ffplayPath == "" {
		ffplayPath = "ffplay"
	}

	et, err := newUtility(ffplayPath)
	if err != nil {
		return nil, err
	}

	return &FFplay{
		externalTool: et,
	}, nil
}

// Path returns the path to the ffplay executable
func (fp *FFplay) Path() string {
	return fp.utilityPath
}

// Play invokes ffplay to play from streaming audio input
func (fp *FFplay) Play(appCtx *synapseq.AppContext) error {
	if appCtx == nil {
		return fmt.Errorf("app context cannot be nil")
	}

	ffplay := exec.Command(
		fp.utilityPath,
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
