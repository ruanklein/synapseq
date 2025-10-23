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

	"github.com/ruanklein/synapseq/internal/audio"
	t "github.com/ruanklein/synapseq/internal/types"
)

// extract handles the extraction of a text sequence from a WAV file
func extract(app *t.AppContext) error {
	if app == nil {
		return nil
	}

	var content string
	var err error

	if content, err = audio.ExtractTextSequenceFromWAV(app.InputFile); err != nil {
		return fmt.Errorf("extract: %v", err)
	}

	if app.OutputFile == "-" {
		fmt.Print(content)
		return nil
	}

	if err = os.WriteFile(app.OutputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("extract: %v", err)
	}

	if !app.Quiet {
		fmt.Fprintf(os.Stderr, "extract: extracted text sequence to %s\n", app.OutputFile)
	}

	return nil
}
