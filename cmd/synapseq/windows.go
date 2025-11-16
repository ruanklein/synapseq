//go:build windows

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package main

import (
	"fmt"

	"github.com/ruanklein/synapseq/v3/internal/fileassoc"
)

// installWindowsFileAssociation sets up the file association for .spsq files on Windows
func installWindowsFileAssociation(quiet bool) error {
	_ = fileassoc.CleanSynapSeqWindowsRegistry()

	if err := fileassoc.InstallWindowsFileAssociation(); err != nil {
		return err
	}
	if err := fileassoc.InstallWindowsContextMenu(); err != nil {
		return err
	}
	if err := fileassoc.InstallWindowsExtractMenu(); err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Successfully installed .spsq file association with SynapSeq.")
	}
	return nil
}

// uninstallWindowsFileAssociation removes the file association for .spsq files on Windows
func uninstallWindowsFileAssociation(quiet bool) error {
	if err := fileassoc.CleanSynapSeqWindowsRegistry(); err != nil {
		return err
	}

	if !quiet {
		fmt.Println("Successfully removed .spsq file association with SynapSeq.")
	}
	return nil
}
