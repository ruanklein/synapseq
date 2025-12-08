//go:build nohub

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
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
