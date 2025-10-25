/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	"os"

	"github.com/ruanklein/synapseq/internal/audio"
)

// extract extracts the text sequence from a WAV file
func (ac *AppContext) extract() (string, error) {
	content, err := audio.ExtractTextSequenceFromWAV(ac.inputFile)
	if err != nil {
		return "", err
	}

	return content, nil
}

// Extract generates the extracted text sequence from the WAV input file
func (ac *AppContext) Extract() (string, error) {
	content, err := ac.extract()
	if err != nil {
		return "", err
	}

	return content, nil
}

// SaveExtracted saves the extracted text sequence to the output file
func (ac *AppContext) SaveExtracted() error {
	content, err := ac.extract()
	if err != nil {
		return err
	}

	if err = os.WriteFile(ac.outputFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
