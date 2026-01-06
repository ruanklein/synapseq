//go:build nohub

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

package main

import (
	"fmt"

	"github.com/ruanklein/synapseq/v3/internal/cli"
)

// hubRunUpdate is disabled when built with -tags=nohub
func hubRunUpdate(quiet bool) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunClean is disabled when built with -tags=nohub
func hubRunClean(quiet bool) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunGet is disabled when built with -tags=nohub
func hubRunGet(sequenceId, outputFile string, opts *cli.CLIOptions) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunList is disabled when built with -tags=nohub
func hubRunList() error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunSearch is disabled when built with -tags=nohub
func hubRunSearch(query string) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunDownload is disabled when built with -tags=nohub
func hubRunDownload(sequenceID, targetDir string, quiet bool) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}

// hubRunInfo is disabled when built with -tags=nohub
func hubRunInfo(sequenceID string) error {
	return fmt.Errorf("Hub functionality is disabled in this build")
}
