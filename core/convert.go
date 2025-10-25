/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	"fmt"
	"os"

	seq "github.com/ruanklein/synapseq/internal/sequence"
	t "github.com/ruanklein/synapseq/internal/types"
)

// convert converts the loaded sequence to text format
func (ac *AppContext) convert() (string, error) {
	if ac.format == t.FormatText {
		return "", fmt.Errorf("input format is already text")
	}

	if ac.sequence == nil {
		return "", fmt.Errorf("sequence is nil")
	}

	content, err := seq.ConvertToText(ac.sequence)
	if err != nil {
		return "", err
	}

	return content, nil
}

// Text generates the text sequence from the loaded sequence
func (ac *AppContext) Text() (string, error) {
	content, err := ac.convert()
	if err != nil {
		return "", err
	}

	return content, nil
}

// SaveText saves the text sequence to the output file
func (ac *AppContext) SaveText() error {
	content, err := ac.convert()
	if err != nil {
		return err
	}

	if err = os.WriteFile(ac.outputFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
