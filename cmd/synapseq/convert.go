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

	"github.com/ruanklein/synapseq/internal/sequence"
	t "github.com/ruanklein/synapseq/internal/types"
)

// convert handles the conversion of a structured sequence to a text-based sequence file
func convert(app *t.AppContext) error {
	if app == nil {
		return nil
	}

	if app.Format == t.FormatText {
		return fmt.Errorf("convert: input format is already text")
	}

	var content string
	var err error

	if content, err = sequence.ConvertToText(app.Sequence); err != nil {
		return fmt.Errorf("convert: %v", err)
	}

	if app.OutputFile == "-" {
		fmt.Print(content)
		return nil
	}

	if err = os.WriteFile(app.OutputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("convert: %v", err)
	}

	if !app.Quiet {
		fmt.Fprintf(os.Stderr, "convert: converted sequence from %s to text format: %s\n", app.Format.String(), app.OutputFile)
	}

	return nil
}
