//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package core

import (
	"os"

	"github.com/ruanklein/synapseq/v3/internal/audio"
)

// extract extracts the text sequence from a WAV file
func extract(inputFile string) (string, error) {
	content, err := audio.ExtractTextSequenceFromWAV(inputFile)
	if err != nil {
		return "", err
	}

	return content, nil
}

// Extract generates the extracted text sequence from the WAV input file
func Extract(inputFile string) (string, error) {
	content, err := extract(inputFile)
	if err != nil {
		return "", err
	}

	return content, nil
}

// SaveExtracted saves the extracted text sequence to the output file
func SaveExtracted(inputFile, outputFile string) error {
	content, err := extract(inputFile)
	if err != nil {
		return err
	}

	if err = os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
