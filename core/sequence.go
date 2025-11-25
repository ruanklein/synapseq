//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	seq "github.com/ruanklein/synapseq/v3/internal/sequence"
	t "github.com/ruanklein/synapseq/v3/internal/types"
)

// LoadSequence loads the sequence from the input file based on the specified format
func (ac *AppContext) LoadSequence() error {
	var err error
	if ac.format == t.FormatText {
		ac.sequence, err = seq.LoadTextSequence(ac.inputFile)
	} else {
		ac.sequence, err = seq.LoadStructuredSequence(ac.inputFile, ac.format)
	}

	if err != nil {
		return err
	}
	return nil
}

// Comments returns the comments from the loaded sequence
func (ac *AppContext) Comments() []string {
	if ac.sequence == nil {
		return nil
	}
	return ac.sequence.Comments
}
