//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
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
